package handler_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	_ "modernc.org/sqlite"
	"quran-api-go/internal/handler"
	"quran-api-go/internal/repository"
	"quran-api-go/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestJuzHandler_ListResponse_IncludesTotalAyahs verifies list response has total_ayahs
func TestJuzHandler_ListResponse_IncludesTotalAyahs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Setup schema - need both tables for the JOIN
	db.ExecContext(context.Background(), `
		CREATE TABLE juzs (id INTEGER PRIMARY KEY, juz_number INTEGER, first_ayah_id INTEGER, last_ayah_id INTEGER);
		CREATE TABLE ayahs (id INTEGER PRIMARY KEY, juz_number INTEGER);

		INSERT INTO juzs (id, juz_number, first_ayah_id, last_ayah_id) VALUES (1, 1, 1, 7);
		-- Insert 7 ayahs for juz 1
		INSERT INTO ayahs (id, juz_number) VALUES (1, 1), (2, 1), (3, 1), (4, 1), (5, 1), (6, 1), (7, 1);
	`)

	repo := repository.NewJuzRepository(db)
	svc := service.NewJuzService(repo)
	h := handler.NewJuzHandler(svc)
	r.GET("/juz", h.List)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/juz", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	data, ok := body["data"].([]any)
	if !ok {
		t.Fatalf("expected data array, got %T", body["data"])
	}

	if len(data) == 0 {
		t.Fatal("expected at least one juz")
	}

	juz := data[0].(map[string]any)
	if juz["total_ayahs"] == nil {
		t.Errorf("Juz list response missing 'total_ayahs' field. Got keys: %v", getMapKeys(juz))
	}

	// Verify the count is correct
	totalAyahs := int(juz["total_ayahs"].(float64))
	if totalAyahs != 7 {
		t.Errorf("Expected total_ayahs=7, got %d", totalAyahs)
	}
}

// TestJuzHandler_AyahsResponse_WrappedStructure verifies ayahs response has correct structure
func TestJuzHandler_AyahsResponse_WrappedStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.ExecContext(context.Background(), `
		CREATE TABLE surahs (id INTEGER PRIMARY KEY, name_latin TEXT);
		CREATE TABLE ayahs (id INTEGER PRIMARY KEY, surah_id INTEGER, number_in_surah INTEGER, text_uthmani TEXT, translation_indo TEXT, translation_en TEXT, juz_number INTEGER);
		CREATE TABLE juzs (id INTEGER PRIMARY KEY, juz_number INTEGER, first_ayah_id INTEGER, last_ayah_id INTEGER, total_ayahs INTEGER);

		INSERT INTO surahs (id, name_latin) VALUES (1, 'Al-Fatihah');
		INSERT INTO juzs (id, juz_number, first_ayah_id, last_ayah_id, total_ayahs) VALUES (1, 1, 1, 7, 7);
		INSERT INTO ayahs (id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number)
		VALUES (1, 1, 1, 'bismillah', 'bismillah', 'bismillah', 1);
	`)

	repo := repository.NewJuzRepository(db)
	svc := service.NewJuzService(repo)
	h := handler.NewJuzHandler(svc)
	r.GET("/juz/:number/ayah", h.Ayahs)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/juz/1/ayah", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %T", body["data"])
	}

	if data["juz"] == nil {
		t.Errorf("Response missing 'juz' wrapper. Got keys: %v", getMapKeys(data))
	}
	if data["ayahs"] == nil {
		t.Errorf("Response missing 'ayahs' array. Got keys: %v", getMapKeys(data))
	}
}

func getMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
