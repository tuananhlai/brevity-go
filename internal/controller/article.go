package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/service"
)

type ArticleController struct {
	articleService service.ArticleService
}

func NewArticleController(articleService service.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

func (c *ArticleController) ListPreviews(ctx *gin.Context) {
	var req ListPreviewsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    CodeBindingRequestError,
			Message: err.Error(),
		})
		return
	}

	req.PageSize = clamp(req.PageSize, 1, 100)

	articles, nextPageToken, err := c.articleService.ListPreviews(ctx.Request.Context(),
		req.PageSize, repository.WithPageToken(req.PageToken))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    CodeUnknown,
			Message: err.Error(),
		})
		return
	}

	response := ListPreviewsResponse{
		Items:         make([]ListPreviewResponseItem, len(articles)),
		NextPageToken: nextPageToken,
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

type ListPreviewsRequest struct {
	PageToken string `form:"pageToken"`
	PageSize  int    `form:"pageSize,default=50"`
}

type ListPreviewsResponse struct {
	Items         []ListPreviewResponseItem `json:"items"`
	NextPageToken string                    `json:"nextPageToken"`
}
