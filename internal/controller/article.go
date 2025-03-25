package controller

import (
	"encoding/json"
	"net/http"

	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/service"
)

type ArticleController struct {
	svc service.ArticleService
}

func NewArticleController(svc service.ArticleService) *ArticleController {
	return &ArticleController{svc: svc}
}

// CreateArticle handles article creation requests
func (c *ArticleController) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article model.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.svc.Create(r.Context(), &article); err != nil {
		http.Error(w, "Failed to create article", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ListArticlePreviews handles requests to list article previews
func (c *ArticleController) ListArticlePreviews(w http.ResponseWriter, r *http.Request) {
	previews, err := c.svc.ListPreviews(r.Context())
	if err != nil {
		http.Error(w, "Failed to list articles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(previews)
}
