package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"

	"github.com/tuananhlai/brevity-go/internal/store"
	"github.com/tuananhlai/brevity-go/internal/token"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

const (
	CodeInvalidCredentials ErrorCode = "invalid_credentials"
	CodeUserAlreadyExists  ErrorCode = "user_already_exists"
)

// AuthStore defines the store methods used by the auth controller.
type AuthStore interface {
	GetUser(ctx context.Context, emailOrUsername string) (*store.User, error)
	CreateUser(ctx context.Context, params store.CreateUserParams) (*store.User, error)
	GetUserByID(ctx context.Context, userID string) (*store.User, error)
}

type AuthController struct {
	store       AuthStore
	tokenIssuer *token.AccessTokenIssuer
}

func NewAuthController(store AuthStore, tokenIssuer *token.AccessTokenIssuer) *AuthController {
	return &AuthController{
		store:       store,
		tokenIssuer: tokenIssuer,
	}
}

func (c *AuthController) Login(ginCtx *gin.Context) {
	ctx, span := otel.Tracer(packageName).Start(ginCtx.Request.Context(), "AuthController.Login")
	defer span.End()

	var req LoginRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		WriteBindingErrorResponse(ginCtx, span, err)
		return
	}

	span.SetAttributes(
		attribute.String("emailOrUsername", req.EmailOrUsername),
	)

	user, err := c.store.GetUser(ctx, req.EmailOrUsername)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			WriteErrorResponse(ginCtx, WriteErrorResponseParams{
				Body: ErrorResponse{
					Code:    CodeInvalidCredentials,
					Message: fmt.Sprintf("%s: %s", ErrInvalidCredentials, err),
				},
				Span: span,
				Err:  ErrInvalidCredentials,
			})
			return
		}

		WriteUnknownErrorResponse(ginCtx, span, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(req.Password)); err != nil {
		WriteErrorResponse(ginCtx, WriteErrorResponseParams{
			Body: ErrorResponse{
				Code:    CodeInvalidCredentials,
				Message: fmt.Sprintf("%s: %s", ErrInvalidCredentials, err),
			},
			Span: span,
			Err:  ErrInvalidCredentials,
		})
		return
	}

	accessToken, err := c.tokenIssuer.Issue(user.ID.String())
	if err != nil {
		WriteUnknownErrorResponse(ginCtx, span, err)
		return
	}

	res := LoginResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}

	setAccessTokenCookie(ginCtx, accessToken)
	ginCtx.JSON(http.StatusOK, res)
}

func (c *AuthController) Register(ginCtx *gin.Context) {
	ctx, span := otel.Tracer(packageName).Start(ginCtx.Request.Context(), "AuthController.Register")
	defer span.End()

	var req RegisterRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		WriteBindingErrorResponse(ginCtx, span, err)
		return
	}

	span.SetAttributes(
		attribute.String("email", req.Email),
		attribute.String("username", req.Username),
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteUnknownErrorResponse(ginCtx, span, err)
		return
	}

	_, err = c.store.CreateUser(ctx, store.CreateUserParams{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		if errors.Is(err, store.ErrUserAlreadyExists) {
			WriteErrorResponse(ginCtx, WriteErrorResponseParams{
				Body: ErrorResponse{
					Code:    CodeUserAlreadyExists,
					Message: fmt.Sprintf("%s: %s", store.ErrUserAlreadyExists, err),
				},
				Span: span,
				Err:  err,
			})
			return
		}

		WriteUnknownErrorResponse(ginCtx, span, err)
		return
	}

	ginCtx.Status(http.StatusOK)
}

func (c *AuthController) GetCurrentUser(ginCtx *gin.Context) {
	ctx, span := otel.Tracer(packageName).Start(ginCtx.Request.Context(), "AuthController.GetCurrentUser")
	defer span.End()

	userID, err := getContextUserID(ginCtx)
	if err != nil {
		WriteErrorResponse(ginCtx, WriteErrorResponseParams{
			Body: ErrorResponse{
				Code:    CodeUnauthorized,
				Message: err.Error(),
			},
			Span: span,
			Err:  err,
		})
		return
	}

	user, err := c.store.GetUserByID(ctx, userID)
	if err != nil {
		WriteUnknownErrorResponse(ginCtx, span, err)
		return
	}

	res := GetCurrentUserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}

	ginCtx.JSON(http.StatusOK, res)
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

type GetCurrentUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
