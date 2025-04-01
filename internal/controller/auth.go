package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuananhlai/brevity-go/internal/service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ginCtx *gin.Context) {
	ctx, span := appTracer.Start(ginCtx.Request.Context(), "AuthController.Register")
	defer span.End()

	var req RegisterRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		span.SetStatus(codes.Error, "failed to bind request")
		span.RecordError(err)
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(
		attribute.String("email", req.Email),
		attribute.String("username", req.Username),
	)

	err := c.authService.Register(ctx, req.Email, req.Username, req.Password)
	if err != nil {
		span.SetStatus(codes.Error, "failed to register user")
		span.RecordError(err)
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
