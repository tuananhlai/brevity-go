package service

import (
	"context"

	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/repository"
)

// ArticleService defines the interface for article business logic
type ArticleService interface {
	Create(ctx context.Context, article *model.Article) error
	ListPreviews(ctx context.Context, pageSize int, opts ...repository.ListPreviewsOption) ([]model.ArticlePreview, string, error)
}

type articleServiceImpl struct {
	repo repository.ArticleRepository
}

// NewArticleService creates a new article service
func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleServiceImpl{repo: repo}
}

// Create creates a new article
func (s *articleServiceImpl) Create(ctx context.Context, article *model.Article) error {
	return s.repo.Create(ctx, article)
}

// ListPreviews lists articles with basic information
func (s *articleServiceImpl) ListPreviews(ctx context.Context,
	pageSize int,
	opts ...repository.ListPreviewsOption,
) ([]model.ArticlePreview, string, error) {
	previews, nextPageToken, err := s.repo.ListPreviews(ctx, pageSize, opts...)
	if err != nil {
		return nil, "", err
	}

	return previews, nextPageToken, nil
}
