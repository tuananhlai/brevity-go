package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/tuananhlai/brevity-go/internal/service"
)

type ArticleController struct {
	articleService service.ArticleService
}

func NewArticleController(articleService service.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

func (c *ArticleController) ListPreviews(ctx *gin.Context) {
	articles, err := c.articleService.ListPreviews(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ListPreviewsResponse{
		Items: make([]ListPreviewResponseItem, len(articles)),
	}
	for i, article := range articles {
		response.Items[i] = ListPreviewResponseItem{
			ID:                article.ID,
			Slug:              article.Slug,
			Title:             article.Title,
			Description:       article.Description,
			AuthorID:          article.AuthorID,
			AuthorDisplayName: article.AuthorDisplayName,
			CreatedAt:         article.CreatedAt,
			UpdatedAt:         article.UpdatedAt,
		}
	}
	ctx.JSON(http.StatusOK, response)
}

type ListPreviewResponseItem struct {
	ID                uuid.UUID `json:"id"`
	Slug              string    `json:"slug"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	AuthorID          uuid.UUID `json:"authorID"`
	AuthorDisplayName string    `json:"authorDisplayName"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type ListPreviewsResponse struct {
	Items []ListPreviewResponseItem `json:"items"`
}
