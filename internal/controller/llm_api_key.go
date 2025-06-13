package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tuananhlai/brevity-go/internal/service"
)

type LLMAPIKeyController struct {
	llmAPIKeyService service.LLMAPIKeyService
}

func NewLLMAPIKeyController(llmAPIKeyService service.LLMAPIKeyService) *LLMAPIKeyController {
	return &LLMAPIKeyController{llmAPIKeyService: llmAPIKeyService}
}

func (c *LLMAPIKeyController) RegisterRoutes(router *gin.Engine) {
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

	ctx, span := appTracer.Start(ginCtx.Request.Context(), "LLMAPIKeyController.CreateLLMAPIKey")
	defer span.End()

	var req request
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	llmAPIKey, err := c.llmAPIKeyService.Create(ctx, service.LLMAPIKeyCreateParams{
		Name:  req.Name,
		Value: req.Value,
		// TODO: add user ID from context
	})
	if err != nil {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
