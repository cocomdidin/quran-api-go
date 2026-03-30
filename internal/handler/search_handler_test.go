package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"quran-api-go/internal/domain/search"
	"quran-api-go/internal/handler"
	"testing"

	"github.com/gin-gonic/gin"
)

// Mock search service
type mockSearchService struct{}

func (m *mockSearchService) Search(ctx context.Context, p search.Params) ([]search.Result, int, error) {
	return []search.Result{
		{ID: 1, TextUthmani: "test"},
	}, 100, nil
}

func TestSearchHandler_ResponseIncludesQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := &mockSearchService{}
	h := handler.NewSearchHandler(svc)
	r.GET("/search", h.Search)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/search?q=test&lang=id&page=1&limit=20", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	// Check response structure
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data field, got %v", body["data"])
	}

	// Verify query field exists
	if data["query"] == nil {
		t.Errorf("Response missing 'query' field. Got keys: %v", getKeys(data))
	}
	if data["query"] != "test" {
		t.Errorf("Expected query='test', got %v", data["query"])
	}

	// Verify other expected fields
	if data["total"] == nil {
		t.Errorf("Response missing 'total' field")
	}
	if data["page"] == nil {
		t.Errorf("Response missing 'page' field")
	}
	if data["limit"] == nil {
		t.Errorf("Response missing 'limit' field")
	}
	if data["results"] == nil {
		t.Errorf("Response missing 'results' field")
	}
}

func getKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
