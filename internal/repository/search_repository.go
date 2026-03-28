package repository

import (
	"context"
	"database/sql"
	"fmt"

	"quran-api-go/internal/domain/search"
)

type searchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) search.SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) Search(ctx context.Context, p search.Params) ([]search.Result, int, error) {
	// Build FTS5 MATCH query with optional filters
	whereClause := "ayahs_fts MATCH ?"
	args := []interface{}{p.Query + "*"} // prefix match for partial terms

	// Add optional filters
	if p.SurahID > 0 {
		whereClause += " AND a.surah_id = ?"
		args = append(args, p.SurahID)
	}
	if p.Juz > 0 {
		whereClause += " AND a.juz_number = ?"
		args = append(args, p.Juz)
	}

	// Count total results
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT ayahs_fts.ayah_id)
		FROM ayahs_fts
		JOIN ayahs a ON a.id = ayahs_fts.ayah_id
		WHERE %s
	`, whereClause)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	// Fetch paginated results with surah info
	limit := p.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (p.Page - 1) * p.Limit
	if offset < 0 {
		offset = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT a.id, a.surah_id, s.name_latin, a.number_in_surah,
			   a.text_uthmani, a.translation_indo, a.translation_en, a.juz_number
		FROM ayahs_fts
		JOIN ayahs a ON a.id = ayahs_fts.ayah_id
		JOIN surahs s ON a.surah_id = s.id
		WHERE %s
		GROUP BY a.id
		ORDER BY a.id ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []search.Result
	for rows.Next() {
		var r search.Result
		var translationIndo, translationEn string

		if err := rows.Scan(
			&r.ID,
			&r.SurahID,
			&r.SurahInfo.NameLatin,
			&r.NumberInSurah,
			&r.TextUthmani,
			&translationIndo,
			&translationEn,
			&r.JuzNumber,
		); err != nil {
			return nil, 0, err
		}

		r.SurahInfo.ID = r.SurahID

		// Set translation based on lang
		if p.Lang == "en" {
			r.Translation = translationEn
		} else {
			r.Translation = translationIndo
		}

		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
