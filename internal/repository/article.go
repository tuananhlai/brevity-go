package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/tuananhlai/brevity-go/internal/model"
)

var ErrArticleNotFound = errors.New("article not found")

// ArticleRepository defines the interface for article data access
type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	ListPreviews(ctx context.Context, pageSize int, opts ...ListPreviewsOption) (
		results []model.ArticlePreview, nextPageToken string, err error)
	GetBySlug(ctx context.Context, slug string) (*model.ArticleDetails, error)
}

type articleRepositoryImpl struct {
	db *sqlx.DB
}

// NewArticleRepository creates a new article repository
func NewArticleRepository(db *sqlx.DB) ArticleRepository {
	return &articleRepositoryImpl{db: db}
}

// Create creates a new article
func (r *articleRepositoryImpl) Create(ctx context.Context, article *model.Article) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO articles (slug, title, description, plaintext_content, 
		content, author_id) VALUES ($1, $2, $3, $4, $5, $6)`,
		article.Slug, article.Title, article.Description,
		article.PlaintextContent, article.Content, article.AuthorID)
	if err != nil {
		return err
	}

	return nil
}

func (r *articleRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*model.ArticleDetails, error) {
	article := model.ArticleDetails{}
	err := r.db.GetContext(ctx, &article, `
		SELECT a.id, a.slug, a.title, a.content, a.author_id, a.created_at, a.updated_at,
			u.username AS author_username, u.display_name AS author_display_name, u.avatar_url AS author_avatar_url
		FROM articles a
		INNER JOIN users u ON a.author_id = u.id
		WHERE a.slug = $1`, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	return &article, nil
}

// ListPreviews lists articles with basic information
func (r *articleRepositoryImpl) ListPreviews(ctx context.Context, pageSize int,
	opts ...ListPreviewsOption,
) ([]model.ArticlePreview, string, error) {
	options := &listPreviewsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var token *ListPreviewsPageToken
	if options.pageToken != "" {
		token = &ListPreviewsPageToken{}
		err := parsePageToken(options.pageToken, token)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse page token: %w", err)
		}
	}

	articles := []model.ArticlePreview{}

	queryBuilder := squirrel.Select(
		"a.id", "a.slug", "a.title", "a.description", "a.author_id", "a.created_at", "a.updated_at",
		"u.username AS author_username", "u.display_name AS author_display_name", "u.avatar_url AS author_avatar_url",
	).
		From("articles a").
		InnerJoin("users u ON a.author_id = u.id").
		OrderBy("a.created_at DESC").
		Limit(uint64(pageSize))

	if token != nil {
		// TODO: Review the logic for next page token.
		queryBuilder = queryBuilder.Where("a.created_at <= ? AND a.id != ?", token.CreatedAt, token.ArticleID)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, "", fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = r.db.SelectContext(ctx, &articles, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to execute SQL query: %w", err)
	}

	var nextPageToken string
	if len(articles) > 0 {
		nextPageToken, err = generatePageToken(ListPreviewsPageToken{
			ArticleID: articles[len(articles)-1].ID.String(),
			CreatedAt: articles[len(articles)-1].CreatedAt,
		})
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate next page token: %w", err)
		}
	}

	return articles, nextPageToken, nil
}

type ListPreviewsPageToken struct {
	ArticleID string    `json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
}

type listPreviewsOptions struct {
	pageToken string
}

type ListPreviewsOption func(*listPreviewsOptions)

func WithPageToken(pageToken string) ListPreviewsOption {
	return func(o *listPreviewsOptions) {
		o.pageToken = pageToken
	}
}
