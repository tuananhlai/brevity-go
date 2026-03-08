package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// CreateArticle creates a new article
func (p *PostgresStore) CreateArticle(ctx context.Context, article *Article) error {
	query, args, err := p.qb.
		Insert("articles").
		Columns("slug", "title", "description", "plaintext_content", "content", "content_format", "author_id").
		Values(article.Slug, article.Title, article.Description, article.PlaintextContent,
			article.Content, "text/markdown", article.AuthorID).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	return err
}

// GetArticleBySlug retrieves a single article by its slug
func (p *PostgresStore) GetArticleBySlug(ctx context.Context, slug string) (*ArticleDetails, error) {
	query, args, err := p.qb.
		Select("a.id", "a.slug", "a.title", "a.content", "a.author_id",
			"a.created_at", "a.updated_at", "da.display_name AS author_display_name").
		From("articles a").
		InnerJoin("digital_authors da ON a.author_id = da.id").
		Where("a.slug = ?", slug).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var article ArticleDetails
	err = p.db.GetContext(ctx, &article, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	return &article, nil
}

// ListArticlesPreviews lists articles with basic information
func (p *PostgresStore) ListArticlesPreviews(ctx context.Context) ([]ArticlePreview, error) {
	articles := []ArticlePreview{}

	query, _, err := p.qb.
		Select("a.id", "a.slug", "a.title", "a.description", "a.author_id",
			"a.created_at", "a.updated_at", "da.display_name AS author_display_name").
		From("articles a").
		InnerJoin("digital_authors da ON a.author_id = da.id").
		OrderBy("a.created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = p.db.SelectContext(ctx, &articles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %w", err)
	}

	return articles, nil
}
