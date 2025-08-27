package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/tuananhlai/brevity-go/internal/model"
)

type DigitalAuthorRepository interface {
	ListByUserID(ctx context.Context, userID string) ([]*model.DigitalAuthor, error)
	Create(ctx context.Context, params DigitalAuthorCreateParams) (*model.DigitalAuthor, error)
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
		SELECT id, owner_id, display_name, system_prompt, default_user_prompt, api_key_id, avatar_url, created_at, updated_at
		FROM digital_authors
		WHERE owner_id = $1
	`

	var digitalAuthors []*model.DigitalAuthor
	if err := d.db.SelectContext(ctx, &digitalAuthors, query, userID); err != nil {
		return nil, err
	}

	return digitalAuthors, nil
}

func (d *digitalAuthorRepositoryImpl) Create(ctx context.Context, params DigitalAuthorCreateParams) (*model.DigitalAuthor, error) {
	query := `
		INSERT INTO digital_authors (owner_id, display_name, system_prompt, default_user_prompt, api_key_id, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, owner_id, display_name, system_prompt, default_user_prompt, api_key_id, avatar_url, created_at, updated_at
	`

	digitalAuthor := &model.DigitalAuthor{}
	err := d.db.GetContext(
		ctx,
		digitalAuthor,
		query,
		params.OwnerID,
		params.DisplayName,
		params.SystemPrompt,
		params.DefaultUserPrompt,
		params.APIKeyID,
		params.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	return digitalAuthor, nil
}

type DigitalAuthorCreateParams struct {
	OwnerID           string
	DisplayName       string
	SystemPrompt      string
	DefaultUserPrompt string
	APIKeyID          string
	AvatarURL         string
}
