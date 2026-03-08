package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/tuananhlai/brevity-go/internal/store"
	"github.com/tuananhlai/brevity-go/internal/telemetry"
)

type DigitalAuthorStore interface {
	ListDigitalAuthors(ctx context.Context) ([]*store.DigitalAuthor, error)
	CreateDigitalAuthor(ctx context.Context, params store.CreateDigitalAuthorParams) (*store.DigitalAuthor, error)
}

type DigitalAuthorController struct {
	store DigitalAuthorStore
}

func NewDigitalAuthorController(store DigitalAuthorStore) *DigitalAuthorController {
	return &DigitalAuthorController{
		store: store,
	}
}

func (c *DigitalAuthorController) ListDigitalAuthors(ginCtx *gin.Context) {
	ctx, span := telemetry.Tracer(otelScopeName).Start(ginCtx.Request.Context(),
		"DigitalAuthorController.ListDigitalAuthors")
	defer span.End()

	digitalAuthors, err := c.store.ListDigitalAuthors(ctx)
	if err != nil {
		writeUnknownErrorResponse(ginCtx, span, err)
		return
	}

	var res ListDigitalAuthorsResponse
	for _, da := range digitalAuthors {
		res.Items = append(res.Items, DigitalAuthor{
			ID:           da.ID,
			DisplayName:  da.DisplayName,
			SystemPrompt: da.SystemPrompt,
			CreatedAt:    da.CreatedAt,
		})
	}

	ginCtx.JSON(http.StatusOK, res)
}

func (c *DigitalAuthorController) CreateDigitalAuthor(ginCtx *gin.Context) {
	ctx, span := telemetry.Tracer(otelScopeName).Start(ginCtx.Request.Context(),
		"DigitalAuthorController.CreateDigitalAuthor")
	defer span.End()

	var req CreateDigitalAuthorRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		writeBindingErrorResponse(ginCtx, span, err)
		return
	}

	da, err := c.store.CreateDigitalAuthor(ctx, store.CreateDigitalAuthorParams{
		DisplayName:  req.DisplayName,
		SystemPrompt: req.SystemPrompt,
	})
	if err != nil {
		writeUnknownErrorResponse(ginCtx, span, err)
		return
	}

	ginCtx.JSON(http.StatusCreated, CreateDigitalAuthorResponse{
		ID:           da.ID,
		DisplayName:  da.DisplayName,
		SystemPrompt: da.SystemPrompt,
		CreatedAt:    da.CreatedAt,
	})
}

type CreateDigitalAuthorRequest struct {
	DisplayName  string `json:"displayName" binding:"required"`
	SystemPrompt string `json:"systemPrompt" binding:"required"`
}

type CreateDigitalAuthorResponse = DigitalAuthor

type ListDigitalAuthorsResponse struct {
	Items []DigitalAuthor `json:"items"`
}

type DigitalAuthor struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"displayName"`
	// TODO: limit length
	SystemPrompt string    `json:"systemPrompt"`
	CreatedAt    time.Time `json:"createdAt"`
}
