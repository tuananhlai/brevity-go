package controller

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	"github.com/tuananhlai/brevity-go/internal/store"
)

const (
	CodeArticleNotFound ErrorCode = "article_not_found"
)

// ArticleStore defines the store methods used by the article controller.
type ArticleStore interface {
	CreateArticle(ctx context.Context, article *store.Article) error
	ListArticlesPreviews(ctx context.Context) ([]store.ArticlePreview, error)
	GetArticleBySlug(ctx context.Context, slug string) (*store.ArticleDetails, error)
}

type ArticleController struct {
	store ArticleStore
}

func NewArticleController(store ArticleStore) *ArticleController {
	return &ArticleController{store: store}
}

func (c *ArticleController) ListPreviews(ginCtx *gin.Context) {
	ctx, span := otel.Tracer(otelScopeName).Start(ginCtx.Request.Context(), "ArticleController.ListPreviews")
	defer span.End()

	articles, err := c.store.ListArticlesPreviews(ctx)
	if err != nil {
		writeUnknownErrorResponse(ginCtx, span, err)
		return
	}

	response := ListPreviewsResponse{
		Items: make([]ArticlePreview, len(articles)),
	}
	for i, article := range articles {
		response.Items[i] = ArticlePreview{
			ID:          article.ID,
			Slug:        article.Slug,
			Title:       article.Title,
			Description: article.Description,
			Author: ArticlePreviewAuthor{
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
	ctx, span := otel.Tracer(otelScopeName).Start(ginCtx.Request.Context(), "ArticleController.GetBySlug")
	defer span.End()

	var req GetBySlugRequest
	if err := ginCtx.ShouldBindUri(&req); err != nil {
		writeBindingErrorResponse(ginCtx, span, err)
		return
	}

	article, err := c.store.GetArticleBySlug(ctx, req.Slug)
	if err != nil {
		if errors.Is(err, store.ErrArticleNotFound) {
			writeErrorResponse(ginCtx, writeErrorResponseParams{
				Body: ErrorResponse{
					Code:    CodeArticleNotFound,
					Message: err.Error(),
				},
				Span: span,
				Err:  err,
			})
			return
		}

		writeUnknownErrorResponse(ginCtx, span, err)
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

type ArticlePreview struct {
	ID          uuid.UUID            `json:"id"`
	Slug        string               `json:"slug"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Author      ArticlePreviewAuthor `json:"author"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
}

type ArticlePreviewAuthor struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName,omitempty"`
	AvatarURL   string    `json:"avatarURL,omitempty"`
}

type ListPreviewsResponse struct {
	Items []ArticlePreview `json:"items"`
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
