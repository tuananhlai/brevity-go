package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/tuananhlai/brevity-go/internal/repository"
	"github.com/tuananhlai/brevity-go/internal/service"
)

const (
	CodeArticleNotFound Code = "article_not_found"
)

type ArticleController struct {
	articleService service.ArticleService
}

func NewArticleController(articleService service.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

func (c *ArticleController) RegisterRoutes(router *gin.Engine) {
	router.GET("/v1/article-previews", c.ListPreviews)
	router.GET("/v1/articles/:slug", c.GetBySlug)
}

func (c *ArticleController) ListPreviews(ginCtx *gin.Context) {
	appLogger.Info("Processing ListPreviews request...")

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
	var req GetBySlugRequest
	if err := ginCtx.ShouldBindUri(&req); err != nil {
		ginCtx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    CodeBindingRequestError,
			Message: err.Error(),
		})
		return
	}

	article, err := c.articleService.GetBySlug(ginCtx.Request.Context(), req.Slug)
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			ginCtx.JSON(http.StatusNotFound, ErrorResponse{
				Code:    CodeArticleNotFound,
				Message: err.Error(),
			})
			return
		}

		ginCtx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    CodeUnknown,
			Message: err.Error(),
		})
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
