package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/tuananhlai/brevity-go/internal/utils"
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

// ListArticlesPreviews lists articles with basic information
func (p *Postgres) ListArticlesPreviews(ctx context.Context, pageSize int,
	opts ...ListArticlesPreviewsOption,
) ([]ArticlePreview, string, error) {
	options := &listArticlesPreviewsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var token *ListArticlesPreviewsPageToken
	if options.pageToken != "" {
		token = &ListArticlesPreviewsPageToken{}
		err := utils.ParsePageToken(options.pageToken, token)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse page token: %w", err)
		}
	}

	articles := []ArticlePreview{}

	queryBuilder := p.qb.Select(
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

	err = p.db.SelectContext(ctx, &articles, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to execute SQL query: %w", err)
	}

	var nextPageToken string
	if len(articles) > 0 {
		nextPageToken, err = utils.GeneratePageToken(ListArticlesPreviewsPageToken{
			ArticleID: articles[len(articles)-1].ID.String(),
			CreatedAt: articles[len(articles)-1].CreatedAt,
		})
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate next page token: %w", err)
		}
	}

	return articles, nextPageToken, nil
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
