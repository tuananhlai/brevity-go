package article

import "context"

type Repository interface {
	Create(ctx context.Context, article *Article) error
	ListPreviews(ctx context.Context) ([]ArticlePreview, error)
}

type Service struct {
	repo Repository
}

func NewService(repo *repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, article *Article) error {
	return s.repo.Create(ctx, article)
}

func (s *Service) ListPreviews(ctx context.Context) ([]ArticlePreview, error) {
	return s.repo.ListPreviews(ctx)
}
