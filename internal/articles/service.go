package articles

import (
	"context"
	"errors"

	"github.com/tuananhlai/brevity-go/internal/repository"
)

var ErrArticleNotFound = errors.New("article not found")

type Service interface {
	Create(ctx context.Context, article *repository.Article) error
	ListPreviews(ctx context.Context, pageSize int, opts ...repository.ListArticlesPreviewsOption) (
		[]repository.ArticlePreview, string, error)
	GetBySlug(ctx context.Context, slug string) (*repository.ArticleDetails, error)
}

type serviceImpl struct {
	repo repository.Repository
}

// NewService creates a new article service
func NewService(repo repository.Repository) *serviceImpl {
	return &serviceImpl{repo: repo}
}

// Create creates a new article
func (s *serviceImpl) Create(ctx context.Context, article *repository.Article) error {
	return s.repo.CreateArticle(ctx, article)
}

// ListPreviews lists articles with basic information
func (s *serviceImpl) ListPreviews(ctx context.Context,
	pageSize int,
	opts ...repository.ListArticlesPreviewsOption,
) ([]repository.ArticlePreview, string, error) {
	previews, nextPageToken, err := s.repo.ListArticlesPreviews(ctx, pageSize, opts...)
	if err != nil {
		return nil, "", err
	}

	return previews, nextPageToken, nil
}

func (s *serviceImpl) GetBySlug(ctx context.Context, slug string) (*repository.ArticleDetails, error) {
	article, err := s.repo.GetArticleBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrArticleNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return article, nil
}
