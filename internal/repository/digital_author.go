package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/tuananhlai/brevity-go/internal/model"
)

type DigitalAuthorRepository interface {
	ListByUserID(ctx context.Context, userID string) ([]*model.DigitalAuthor, error)
}

type digitalAuthorRepositoryImpl struct {
	db *sqlx.DB
}

func NewDigitalAuthorRepository(db *sqlx.DB) DigitalAuthorRepository {
	return &digitalAuthorRepositoryImpl{
		db: db,
	}
}

// ListByUserID implements DigitalAuthorRepository.
func (d *digitalAuthorRepositoryImpl) ListByUserID(ctx context.Context, userID string) ([]*model.DigitalAuthor, error) {
	query := `
		SELECT id, owner_id, display_name, system_prompt, avatar_url, created_at, updated_at
		FROM digital_authors
		WHERE owner_id = $1
	`

	var digitalAuthors []*model.DigitalAuthor
	if err := d.db.SelectContext(ctx, &digitalAuthors, query, userID); err != nil {
		return nil, err
	}

	return digitalAuthors, nil
}

type DigitalAuthorCreateParams struct {
	OwnerID      string
	DisplayName  string
	SystemPrompt string
	AvatarURL    string
}
