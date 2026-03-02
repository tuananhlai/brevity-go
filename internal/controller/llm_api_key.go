package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/tuananhlai/brevity-go/internal/llmapikey"
)

type LLMAPIKeyController struct {
	llmAPIKeyManager *llmapikey.Manager
}

func NewLLMAPIKeyController(llmAPIKeyManager *llmapikey.Manager) *LLMAPIKeyController {
	return &LLMAPIKeyController{llmAPIKeyManager: llmAPIKeyManager}
}

func (c *LLMAPIKeyController) ListLLMAPIKeys(ginCtx *gin.Context) {
	type responseItem struct {
		ID            string    `json:"id"`
		Name          string    `json:"name"`
		ValueFirstTen string    `json:"valueFirstTen"`
		ValueLastSix  string    `json:"valueLastSix"`
		CreatedAt     time.Time `json:"createdAt"`
	}

	type response struct {
		Items []responseItem `json:"items"`
	}

	ctx, span := otel.Tracer(packageName).Start(ginCtx.Request.Context(), "LLMAPIKeyController.ListLLMAPIKeys")
	defer span.End()

	userID, err := GetContextUserID(ginCtx)
	if err != nil {
		WriteErrorResponse(ginCtx, WriteErrorResponseParams{
			Body: ErrorResponse{
				Code:    CodeUnauthorized,
				Message: "error getting userID from context",
			},
			Span: span,
			Err:  err,
		})
		return
	}
	span.SetAttributes(attribute.String("userID", userID))

	llmAPIKeys, err := c.llmAPIKeyManager.ListByUserID(ctx, userID)
	if err != nil {
		WriteUnknownErrorResponse(ginCtx, span, err)
		return
	}

	var res response
	for _, llmAPIKey := range llmAPIKeys {
		res.Items = append(res.Items, responseItem{
			ID:            llmAPIKey.ID.String(),
			Name:          llmAPIKey.Name,
			ValueFirstTen: llmAPIKey.ValueFirstTen,
			ValueLastSix:  llmAPIKey.ValueLastSix,
			CreatedAt:     llmAPIKey.CreatedAt,
		})
	}

	ginCtx.JSON(http.StatusOK, res)
}

func (c *LLMAPIKeyController) CreateLLMAPIKey(ginCtx *gin.Context) {
	type request struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	type response struct {
		ID            string    `json:"id"`
		Name          string    `json:"name"`
		ValueFirstTen string    `json:"valueFirstTen"`
		ValueLastSix  string    `json:"valueLastSix"`
		CreatedAt     time.Time `json:"createdAt"`
	}

	ctx, span := otel.Tracer(packageName).Start(ginCtx.Request.Context(), "LLMAPIKeyController.CreateLLMAPIKey")
	defer span.End()

	userID, err := GetContextUserID(ginCtx)
	if err != nil {
		WriteErrorResponse(ginCtx, WriteErrorResponseParams{
			Body: ErrorResponse{
				Code:    CodeUnauthorized,
				Message: "error getting userID from context",
			},
			Span: span,
			Err:  err,
		})
		return
	}

	var req request
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		WriteBindingErrorResponse(ginCtx, span, err)
		return
	}

	llmAPIKey, err := c.llmAPIKeyManager.Create(ctx, llmapikey.CreateInput{
		Name:   req.Name,
		Value:  req.Value,
		UserID: userID,
	})
	if err != nil {
		WriteUnknownErrorResponse(ginCtx, span,
			fmt.Errorf("failed to create llm api key: %w", err))
		return
	}

	ginCtx.JSON(http.StatusOK, response{
		ID:            llmAPIKey.ID.String(),
		Name:          llmAPIKey.Name,
		ValueFirstTen: llmAPIKey.ValueFirstTen,
		ValueLastSix:  llmAPIKey.ValueLastSix,
		CreatedAt:     llmAPIKey.CreatedAt,
	})
}
