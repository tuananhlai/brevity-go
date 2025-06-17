package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tuananhlai/brevity-go/internal/controller/shared"
	"github.com/tuananhlai/brevity-go/internal/service"
)

type LLMAPIKeyController struct {
	llmAPIKeyService service.LLMAPIKeyService
}

func NewLLMAPIKeyController(llmAPIKeyService service.LLMAPIKeyService) *LLMAPIKeyController {
	return &LLMAPIKeyController{llmAPIKeyService: llmAPIKeyService}
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

	userID, err := shared.GetContextUserID(ginCtx)
	if err != nil {
		ginCtx.JSON(http.StatusUnauthorized, shared.ErrorResponse{
			Code:    shared.CodeUnauthorized,
			Message: "error getting userID from context",
		})
		return
	}

	ctx, span := appTracer.Start(ginCtx.Request.Context(), "LLMAPIKeyController.CreateLLMAPIKey")
	defer span.End()

	var req request
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		ginCtx.JSON(http.StatusBadRequest, shared.ErrorResponse{
			Code:    shared.CodeBindingRequestError,
			Message: fmt.Sprintf("error binding request: %s", err.Error()),
		})
		return
	}

	llmAPIKey, err := c.llmAPIKeyService.Create(ctx, service.LLMAPIKeyCreateParams{
		Name:   req.Name,
		Value:  req.Value,
		UserID: userID,
	})
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, shared.ErrorResponse{
			Code:    shared.CodeUnknown,
			Message: fmt.Sprintf("error creating llm api key: %s", err.Error()),
		})
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
