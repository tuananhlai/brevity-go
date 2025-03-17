package article

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(article *Article) error {
	_, err := r.db.Exec(
		"INSERT INTO articles (slug, title, description, text_content, content, author_id) VALUES ($1, $2, $3, $4, $5, $6)",
		article.Slug, article.Title, article.Description,
		article.PlaintextContent, article.Content, article.AuthorID)

	return err
}
