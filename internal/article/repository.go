package article

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, article *Article) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO articles (slug, title, description, plaintext_content, 
		content, author_id) VALUES ($1, $2, $3, $4, $5, $6)`,
		article.Slug, article.Title, article.Description,
		article.PlaintextContent, article.Content, article.AuthorID)

	return err
}

func (r *Repository) ListPreviews(ctx context.Context) ([]ArticlePreview, error) {
	articles := []ArticlePreview{}

	err := r.db.SelectContext(ctx, &articles,
		`SELECT a.id, a.slug, a.title, a.description, a.author_id, a.created_at, a.updated_at, u.display_name FROM articles a INNER JOIN users u ON a.author_id = u.id`,
	)

	return articles, err
}

type ArticlePreview struct {
	ID                uuid.UUID `db:"id"`
	Slug              string    `db:"slug"`
	Title             string    `db:"title"`
	Description       string    `db:"description"`
	AuthorID          uuid.UUID `db:"author_id"`
	AuthorDisplayName string    `db:"display_name"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}
