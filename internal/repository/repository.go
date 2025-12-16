package repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrArticleNotFound   = errors.New("article not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Repository interface {
	CreateArticle(ctx context.Context, article *Article) error
	ListArticlesPreviews(ctx context.Context, pageSize int, opts ...ListArticlesPreviewsOption) (
		results []ArticlePreview, nextPageToken string, err error)
	GetArticleBySlug(ctx context.Context, slug string) (*ArticleDetails, error)
	// GetUser returns the user with the given email or username.
	GetUser(ctx context.Context, emailOrUsername string) (*User, error)
	// CreateUser creates a new user and returns the created user.
	CreateUser(ctx context.Context, params CreateUserParams) (*User, error)
	// GetUserByID returns the user with the given ID.
	GetUserByID(ctx context.Context, userID string) (*User, error)
	ListLLMAPIKeysByUserID(ctx context.Context, userID string) ([]*StoredAPIKey, error)
	CreateLLMAPIKey(ctx context.Context, apiKey CreateLLMAPIKeyParams) (*StoredAPIKey, error)
}

type CreateUserParams struct {
	Email        string
	PasswordHash []byte
	Username     string
}

type CreateLLMAPIKeyParams struct {
	Name         string
	EncryptedKey []byte
	UserID       string
}

type ListArticlesPreviewsPageToken struct {
	ArticleID string    `json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
}

type listArticlesPreviewsOptions struct {
	pageToken string
}

type ListArticlesPreviewsOption func(*listArticlesPreviewsOptions)

func WithPageToken(pageToken string) ListArticlesPreviewsOption {
	return func(o *listArticlesPreviewsOptions) {
		o.pageToken = pageToken
	}
}
