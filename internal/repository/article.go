package repository

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/tuananhlai/brevity-go/internal/model"
)

// ArticleRepository defines the interface for article data access
type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	ListPreviews(ctx context.Context) ([]model.ArticlePreview, error)
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
func (r *articleRepositoryImpl) ListPreviews(ctx context.Context) ([]model.ArticlePreview, error) {
	articles := []model.ArticlePreview{}

	err := r.db.SelectContext(ctx, &articles,
		`SELECT a.id, a.slug, a.title, a.description, a.author_id, a.created_at, a.updated_at, 
		CASE WHEN u.display_name IS NOT NULL THEN u.display_name ELSE u.username END AS display_name 
		FROM articles a 
		INNER JOIN users u ON a.author_id = u.id`)
	if err != nil {
		return nil, err
	}

	return articles, err
}
