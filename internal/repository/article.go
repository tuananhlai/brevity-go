package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// CreateArticle creates a new article
func (r *Postgres) CreateArticle(ctx context.Context, article *Article) error {
	// TODO: Add support for other content formats
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO articles (slug, title, description, plaintext_content,
		content, content_format, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		article.Slug, article.Title, article.Description,
		article.PlaintextContent, article.Content, "text/markdown", article.AuthorID)

	return err
}

func (r *Postgres) GetArticleBySlug(ctx context.Context, slug string) (*ArticleDetails, error) {
	article := ArticleDetails{}
	err := r.db.GetContext(ctx, &article, `
		SELECT a.id, a.slug, a.title, a.content, a.author_id, a.created_at, a.updated_at,
			da.display_name AS author_display_name		
		FROM articles a
		INNER JOIN digital_authors da ON a.author_id = da.id
		WHERE a.slug = $1`, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	return &article, nil
}

// ListArticlesPreviews lists articles with basic information
func (p *Postgres) ListArticlesPreviews(ctx context.Context) ([]ArticlePreview, error) {
	articles := []ArticlePreview{}

	queryBuilder := p.qb.Select("a.id", "a.slug", "a.title", "a.description", "a.author_id",
		"a.created_at", "a.updated_at", "da.display_name AS author_display_name").
		From("articles a").
		InnerJoin("digital_authors da ON a.author_id = da.id").
		OrderBy("a.created_at DESC")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = p.db.SelectContext(ctx, &articles, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %w", err)
	}

	return articles, nil
}
