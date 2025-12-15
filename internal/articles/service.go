package articles

import (
	"context"
)

type Service interface {
	Create(ctx context.Context, article *Article) error
	ListPreviews(ctx context.Context, pageSize int, opts ...ListPreviewsOption) ([]ArticlePreview, string, error)
	GetBySlug(ctx context.Context, slug string) (*ArticleDetails, error)
}

type serviceImpl struct {
	repo Repository
}

// NewService creates a new article service
func NewService(repo Repository) *serviceImpl {
	return &serviceImpl{repo: repo}
}

// Create creates a new article
func (s *serviceImpl) Create(ctx context.Context, article *Article) error {
	return s.repo.Create(ctx, article)
}

// ListPreviews lists articles with basic information
func (s *serviceImpl) ListPreviews(ctx context.Context,
	pageSize int,
	opts ...ListPreviewsOption,
) ([]ArticlePreview, string, error) {
	previews, nextPageToken, err := s.repo.ListPreviews(ctx, pageSize, opts...)
	if err != nil {
		return nil, "", err
	}

	return previews, nextPageToken, nil
}

func (s *serviceImpl) GetBySlug(ctx context.Context, slug string) (*ArticleDetails, error) {
	article, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return article, nil
}
