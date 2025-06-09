package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuananhlai/brevity-go/internal/service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	CodeInvalidCredentials Code = "invalid_credentials"
	CodeUserAlreadyExists  Code = "user_already_exists"

	AccessTokenCookieName   = "access_token"
	AccessTokenCookiePath   = "/"
	AccessTokenCookieMaxAge = 60 * 60 * 24 * 30 // 30 days
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) RegisterRoutes(router *gin.Engine) {
	router.POST("/v1/auth/sign-up", c.Register)
	router.POST("/v1/auth/sign-in", c.Login)
}

func (c *AuthController) Login(ginCtx *gin.Context) {
	ctx, span := appTracer.Start(ginCtx.Request.Context(), "AuthController.Login")
	defer span.End()

	var req LoginRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		span.SetStatus(codes.Error, "failed to bind request")
		span.RecordError(err)
		ginCtx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    CodeBindingRequestError,
			Message: err.Error(),
		})
		return
	}

	span.SetAttributes(
		attribute.String("emailOrUsername", req.EmailOrUsername),
	)

	user, err := c.authService.Login(ctx, req.EmailOrUsername, req.Password)
	if err != nil {
		span.SetStatus(codes.Error, "failed to login")
		span.RecordError(err)

		if errors.Is(err, service.ErrInvalidCredentials) {
			ginCtx.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    CodeInvalidCredentials,
				Message: err.Error(),
			})
			return
		}
		ginCtx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    CodeUnknown,
			Message: err.Error(),
		})
		return
	}

	res := LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	// TODO: See what value to be set to the domain.
	ginCtx.SetCookie(AccessTokenCookieName, user.AccessToken, AccessTokenCookieMaxAge,
		AccessTokenCookiePath, "", true, true)
	ginCtx.JSON(http.StatusOK, res)
}

func (c *AuthController) Register(ginCtx *gin.Context) {
	ctx, span := appTracer.Start(ginCtx.Request.Context(), "AuthController.Register")
	defer span.End()

	var req RegisterRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		span.SetStatus(codes.Error, "failed to bind request")
		span.RecordError(err)
		ginCtx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    CodeBindingRequestError,
			Message: err.Error(),
		})
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

		if errors.Is(err, service.ErrUserAlreadyExists) {
			ginCtx.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    CodeUserAlreadyExists,
				Message: err.Error(),
			})
			return
		}
		ginCtx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    CodeUnknown,
			Message: err.Error(),
		})
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
