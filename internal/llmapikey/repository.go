package llmapikey

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Repository defines database access for LLM API keys.
type Repository interface {
	ListByUserID(ctx context.Context, userID string) ([]*StoredAPIKey, error)
	Create(ctx context.Context, apiKey CreateParams) (*StoredAPIKey, error)
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{
		db: db,
	}
}

func (r *repositoryImpl) ListByUserID(ctx context.Context, userID string) ([]*StoredAPIKey, error) {
	var apiKeys []*StoredAPIKey
	err := r.db.SelectContext(ctx, &apiKeys, `
		SELECT id, name, encrypted_key, created_at, user_id
		FROM llm_api_keys
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}

func (r *repositoryImpl) Create(ctx context.Context, params CreateParams) (*StoredAPIKey, error) {
	apiKey := &StoredAPIKey{}

	err := r.db.GetContext(ctx, apiKey, `
		INSERT INTO llm_api_keys (name, encrypted_key, user_id)
		VALUES ($1, $2, $3) 
		RETURNING id, name, encrypted_key, created_at, user_id`,
		params.Name, params.EncryptedKey, params.UserID)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

type CreateParams struct {
	Name         string
	EncryptedKey []byte
	UserID       string
}
