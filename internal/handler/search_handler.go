package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain/search"
	"quran-api-go/pkg/response"
	"quran-api-go/pkg/validator"
)

type SearchHandler struct {
	service search.SearchService
}

type SearchResponse struct {
	Query   string          `json:"query"`
	Results []search.Result `json:"results"`
	Total   int             `json:"total"`
	Page    int             `json:"page"`
	Limit   int             `json:"limit"`
}

func NewSearchHandler(service search.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "query parameter 'q' is required")
		return
	}

	lang, err := validator.ValidateLang(c.Query("lang"))
	if err != nil {
		response.BadRequest(c, "lang must be 'id' or 'en'")
		return
	}

	surahID, _ := strconv.Atoi(c.Query("surah_id"))
	juz, _ := strconv.Atoi(c.Query("juz"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	params := search.Params{
		Query:   query,
		Lang:    lang,
		SurahID: surahID,
		Juz:     juz,
		Page:    page,
		Limit:   limit,
	}

	results, total, err := h.service.Search(c.Request.Context(), params)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, SearchResponse{
		Query:   query,
		Results: results,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}
