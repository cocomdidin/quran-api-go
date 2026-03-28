package service

import (
	"context"

	"quran-api-go/internal/domain/search"
)

type searchService struct {
	repo search.SearchRepository
}

func NewSearchService(repo search.SearchRepository) search.SearchService {
	return &searchService{repo: repo}
}

func (s *searchService) Search(ctx context.Context, p search.Params) ([]search.Result, int, error) {
	// Set defaults
	if p.Query == "" {
		return []search.Result{}, 0, nil
	}

	if p.Lang == "" {
		p.Lang = "id"
	}
	if p.Lang != "id" && p.Lang != "en" {
		p.Lang = "id"
	}

	if p.Page < 1 {
		p.Page = 1
	}

	if p.Limit < 1 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}

	return s.repo.Search(ctx, p)
}
