package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tuananhlai/brevity-go/internal/model"
)

type LLMAPIKeyRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.LLMAPIKey, error)
	Create(ctx context.Context, apiKey LLMAPIKeyCreateParams) (*model.LLMAPIKey, error)
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

func (r *llmAPIKeyRepositoryImpl) Create(ctx context.Context, params LLMAPIKeyCreateParams) (*model.LLMAPIKey, error) {
	apiKey := &model.LLMAPIKey{}

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

type LLMAPIKeyCreateParams struct {
	Name         string
	EncryptedKey []byte
	UserID       string
}
