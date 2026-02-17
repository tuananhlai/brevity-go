package repository

import (
	"context"
	"errors"
)

var (
	ErrArticleNotFound   = errors.New("article not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Repository interface {
	CreateArticle(ctx context.Context, article *Article) error
	// ListArticlesPreviews returns a list of article previews, which are essentially articles without content.
	ListArticlesPreviews(ctx context.Context, pageSize int, opts ...ListArticlesPreviewsOption) (
		results []ArticlePreview, nextPageToken string, err error)
	GetArticleBySlug(ctx context.Context, slug string) (*ArticleDetails, error)

	// GetUser returns the user with the given email or username.
	GetUser(ctx context.Context, emailOrUsername string) (*User, error)
	// CreateUser creates a new user and returns the created user.
	CreateUser(ctx context.Context, params CreateUserParams) (*User, error)
	// GetUserByID returns the user with the given ID.
	GetUserByID(ctx context.Context, userID string) (*User, error)

	// ListLLMAPIKeysByUserID returns the list of LLM API keys belonging to the given user.
	ListLLMAPIKeysByUserID(ctx context.Context, userID string) ([]*OpenRouterAPIKey, error)
	CreateLLMAPIKey(ctx context.Context, apiKey CreateLLMAPIKeyParams) (*OpenRouterAPIKey, error)

	ListAllDigitalAuthors(ctx context.Context) ([]*DigitalAuthor, error)
}
