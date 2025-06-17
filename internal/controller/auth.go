package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuananhlai/brevity-go/internal/controller/shared"
	"github.com/tuananhlai/brevity-go/internal/service"
	"go.opentelemetry.io/otel/attribute"
)

const (
	CodeInvalidCredentials shared.Code = "invalid_credentials"
	CodeUserAlreadyExists  shared.Code = "user_already_exists"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Login(ginCtx *gin.Context) {
	ctx, span := appTracer.Start(ginCtx.Request.Context(), "AuthController.Login")
	defer span.End()

	var req LoginRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		shared.WriteBindingErrorResponse(ginCtx, span, err)
		return
	}

	span.SetAttributes(
		attribute.String("emailOrUsername", req.EmailOrUsername),
	)

	user, err := c.authService.Login(ctx, req.EmailOrUsername, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			shared.WriteErrorResponse(ginCtx, shared.WriteErrorResponseParams{
				Body: shared.ErrorResponse{
					Code:    CodeInvalidCredentials,
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

	res := LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	shared.SetAccessTokenCookie(ginCtx, user.AccessToken)
	ginCtx.JSON(http.StatusOK, res)
}

func (c *AuthController) Register(ginCtx *gin.Context) {
	ctx, span := appTracer.Start(ginCtx.Request.Context(), "AuthController.Register")
	defer span.End()

	var req RegisterRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		shared.WriteBindingErrorResponse(ginCtx, span, err)
		return
	}

	span.SetAttributes(
		attribute.String("email", req.Email),
		attribute.String("username", req.Username),
	)

	err := c.authService.Register(ctx, req.Email, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			shared.WriteErrorResponse(ginCtx, shared.WriteErrorResponseParams{
				Body: shared.ErrorResponse{
					Code:    CodeUserAlreadyExists,
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

	ginCtx.Status(http.StatusOK)
}

type LoginRequest struct {
	EmailOrUsername string `json:"emailOrUsername" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
