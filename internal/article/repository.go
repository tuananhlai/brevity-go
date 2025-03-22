package article

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

// Create creates a new article.
func (r *repository) Create(ctx context.Context, article *Article) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO articles (slug, title, description, plaintext_content, 
		content, author_id) VALUES ($1, $2, $3, $4, $5, $6)`,
		article.Slug, article.Title, article.Description,
		article.PlaintextContent, article.Content, article.AuthorID)

	return err
}

// ListPreviews lists articles with basic information.
func (r *repository) ListPreviews(ctx context.Context) ([]ArticlePreview, error) {
	articles := []ArticlePreview{}

	err := r.db.SelectContext(ctx, &articles,
		`SELECT a.id, a.slug, a.title, a.description, a.author_id, a.created_at, a.updated_at, CASE WHEN u.display_name IS NOT NULL THEN u.display_name ELSE u.username END AS display_name FROM articles a INNER JOIN users u ON a.author_id = u.id`,
	)

	return articles, err
}
