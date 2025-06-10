package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tuananhlai/brevity-go/internal/model"
)

type LLMAPIKeyRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.LLMAPIKey, error)
	Create(ctx context.Context, apiKey *model.LLMAPIKey) error
}

type llmAPIKeyRepositoryImpl struct {
	db *sqlx.DB
}

func NewLLMAPIKeyRepository(db *sqlx.DB) LLMAPIKeyRepository {
	return &llmAPIKeyRepositoryImpl{
		db: db,
	}
}

func (r *llmAPIKeyRepositoryImpl) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.LLMAPIKey, error) {
	var apiKeys []*model.LLMAPIKey
	err := r.db.SelectContext(ctx, &apiKeys, "SELECT id, name, encrypted_key, created_at, user_id FROM llm_api_keys WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}

func (r *llmAPIKeyRepositoryImpl) Create(ctx context.Context, apiKey *model.LLMAPIKey) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO llm_api_keys (id, name, encrypted_key, user_id) VALUES ($1, $2, $3, $4)",
		apiKey.ID, apiKey.Name, apiKey.EncryptedKey, apiKey.UserID)

	return err
}

// type LLMAPIKey struct {
// 	ID   uuid.UUID
// 	Name string
// 	// Value the plaintext API key.
// 	Value     string
// 	UserID    uuid.UUID
// 	CreatedAt time.Time
// }
