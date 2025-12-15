package controller

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"

	"github.com/tuananhlai/brevity-go/internal/articles"
	"github.com/tuananhlai/brevity-go/internal/controller/shared"
)

const (
	CodeArticleNotFound shared.Code = "article_not_found"
)

// ArticleService defines the interface for article business logic
type ArticleService interface {
	Create(ctx context.Context, article *articles.Article) error
	ListPreviews(ctx context.Context, pageSize int, opts ...articles.ListPreviewsOption) ([]articles.ArticlePreview, string, error)
	GetBySlug(ctx context.Context, slug string) (*articles.ArticleDetails, error)
}

type ArticleController struct {
	articleService ArticleService
}

func NewArticleController(articleService ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

func (c *ArticleController) ListPreviews(ginCtx *gin.Context) {
	appLogger.Info("Processing ListPreviews request...")

	ctx, span := appTracer.Start(ginCtx.Request.Context(), "ArticleController.ListPreviews")
	defer span.End()

	var req ListPreviewsRequest
	if err := ginCtx.ShouldBindQuery(&req); err != nil {
		shared.WriteBindingErrorResponse(ginCtx, span, err)
		return
	}

	span.SetAttributes(attribute.Int("page_size", req.PageSize),
		attribute.String("page_token", req.PageToken))

	req.PageSize = lo.Clamp(req.PageSize, 1, 100)

	articles, nextPageToken, err := c.articleService.ListPreviews(ctx,
		req.PageSize, articles.WithPageToken(req.PageToken))
	if err != nil {
		shared.WriteUnknownErrorResponse(ginCtx, span, err)
		return
	}

	response := ListPreviewsResponse{
		Items:         make([]ListPreviewResponseItem, len(articles)),
		NextPageToken: nextPageToken,
	}
	for i, article := range articles {
		response.Items[i] = ListPreviewResponseItem{
			ID:          article.ID,
			Slug:        article.Slug,
			Title:       article.Title,
			Description: article.Description,
			Author: ListPreviewResponseItemAuthor{
				ID:          article.AuthorID,
				Username:    article.AuthorUsername,
				DisplayName: article.AuthorDisplayName.String,
				AvatarURL:   article.AuthorAvatarURL.String,
			},
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		}
	}
	ginCtx.JSON(http.StatusOK, response)
}

func (c *ArticleController) GetBySlug(ginCtx *gin.Context) {
	ctx, span := appTracer.Start(ginCtx.Request.Context(), "ArticleController.GetBySlug")
	defer span.End()

	var req GetBySlugRequest
	if err := ginCtx.ShouldBindUri(&req); err != nil {
		shared.WriteBindingErrorResponse(ginCtx, span, err)
		return
	}

	article, err := c.articleService.GetBySlug(ctx, req.Slug)
	if err != nil {
		if errors.Is(err, articles.ErrArticleNotFound) {
			shared.WriteErrorResponse(ginCtx, shared.WriteErrorResponseParams{
				Body: shared.ErrorResponse{
					Code:    CodeArticleNotFound,
					Message: err.Error(),
				},
				Span: span,
				Err:  err,
			})
			return
		}

		shared.WriteUnknownErrorResponse(ginCtx, span, err)
		return
	}

	response := GetBySlugResponse{
		ID:      article.ID,
		Slug:    article.Slug,
		Title:   article.Title,
		Content: article.Content,
		Author: GetBySlugResponseAuthor{
			ID:          article.AuthorID,
			Username:    article.AuthorUsername,
			DisplayName: article.AuthorDisplayName.String,
			AvatarURL:   article.AuthorAvatarURL.String,
		},
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}
	ginCtx.JSON(http.StatusOK, response)
}

type ListPreviewResponseItem struct {
	ID          uuid.UUID                     `json:"id"`
	Slug        string                        `json:"slug"`
	Title       string                        `json:"title"`
	Description string                        `json:"description"`
	Author      ListPreviewResponseItemAuthor `json:"author"`
	CreatedAt   time.Time                     `json:"createdAt"`
	UpdatedAt   time.Time                     `json:"updatedAt"`
}

type ListPreviewResponseItemAuthor struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName,omitempty"`
	AvatarURL   string    `json:"avatarURL,omitempty"`
}

type ListPreviewsRequest struct {
	PageToken string `form:"pageToken"`
	PageSize  int    `form:"pageSize,default=50"`
}

type ListPreviewsResponse struct {
	Items         []ListPreviewResponseItem `json:"items"`
	NextPageToken string                    `json:"nextPageToken,omitempty"`
}

type GetBySlugRequest struct {
	Slug string `uri:"slug"`
}

type GetBySlugResponse struct {
	ID        uuid.UUID               `json:"id"`
	Slug      string                  `json:"slug"`
	Title     string                  `json:"title"`
	Content   string                  `json:"content"`
	Author    GetBySlugResponseAuthor `json:"author"`
	CreatedAt time.Time               `json:"createdAt"`
	UpdatedAt time.Time               `json:"updatedAt"`
}

type GetBySlugResponseAuthor struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName,omitempty"`
	AvatarURL   string    `json:"avatarURL,omitempty"`
}
