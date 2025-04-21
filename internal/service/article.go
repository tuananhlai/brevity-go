package service

import (
	"context"
	"errors"

	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/repository"
)

var ErrArticleNotFound = errors.New("article not found")

// ArticleService defines the interface for article business logic
type ArticleService interface {
	Create(ctx context.Context, article *model.Article) error
	ListPreviews(ctx context.Context, pageSize int, opts ...repository.ListPreviewsOption) ([]model.ArticlePreview, string, error)
	GetBySlug(ctx context.Context, slug string) (*model.ArticleDetails, error)
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

func (s *articleServiceImpl) GetBySlug(ctx context.Context, slug string) (*model.ArticleDetails, error) {
	article, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrArticleNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return article, nil
}
