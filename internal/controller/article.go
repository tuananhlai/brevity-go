package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/service"
)

type ArticleController struct {
	articleService service.ArticleService
}

func NewArticleController(articleService service.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

func (c *ArticleController) ListPreviews(ginCtx *gin.Context) {
	ctx, span := appTracer.Start(ginCtx.Request.Context(), "ArticleController.ListPreviews")
	defer span.End()

	var req ListPreviewsRequest
	if err := ginCtx.ShouldBindQuery(&req); err != nil {
		span.SetStatus(codes.Error, "failed to bind request")
		span.RecordError(err)

		ginCtx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    CodeBindingRequestError,
			Message: err.Error(),
		})
		return
	}

	span.SetAttributes(attribute.Int("page_size", req.PageSize),
		attribute.String("page_token", req.PageToken))

	req.PageSize = clamp(req.PageSize, 1, 100)

	articles, nextPageToken, err := c.articleService.ListPreviews(ctx,
		req.PageSize, repository.WithPageToken(req.PageToken))
	if err != nil {
		span.SetStatus(codes.Error, "failed to list previews")
		span.RecordError(err)

		ginCtx.JSON(http.StatusInternalServerError, ErrorResponse{
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
	ginCtx.JSON(http.StatusOK, response)
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
