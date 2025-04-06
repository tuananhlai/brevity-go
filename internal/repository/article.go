package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/tuananhlai/brevity-go/internal/model"
)

// ArticleRepository defines the interface for article data access
type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	ListPreviews(ctx context.Context, pageSize int, opts ...ListPreviewsOption) (
		results []model.ArticlePreview, nextPageToken string, err error)
}

type articleRepositoryImpl struct {
	db     *sqlx.DB
	tracer trace.Tracer
	logger *slog.Logger
}

// NewArticleRepository creates a new article repository
func NewArticleRepository(db *sqlx.DB) ArticleRepository {
	name := "repository.ArticleRepository"

	return &articleRepositoryImpl{
		db:     db,
		tracer: otel.Tracer(name),
		logger: otelslog.NewLogger(name),
	}
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
		"CASE WHEN u.display_name IS NOT NULL THEN u.display_name ELSE u.username END AS display_name",
	).
		From("articles a").
		InnerJoin("users u ON a.author_id = u.id").
		OrderBy("a.created_at DESC").
		Limit(uint64(pageSize))

	if token != nil {
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
