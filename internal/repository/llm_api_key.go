package repository

import (
	"context"
)

func (r *Postgres) ListLLMAPIKeysByUserID(ctx context.Context, userID string) ([]*OpenRouterAPIKey, error) {
	var apiKeys []*OpenRouterAPIKey
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

func (r *Postgres) CreateLLMAPIKey(ctx context.Context, apiKey CreateLLMAPIKeyParams) (*OpenRouterAPIKey, error) {
	storedKey := &OpenRouterAPIKey{}

	err := r.db.GetContext(ctx, storedKey, `
		INSERT INTO llm_api_keys (name, encrypted_key, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, encrypted_key, created_at, user_id`,
		apiKey.Name, apiKey.EncryptedKey, apiKey.UserID)
	if err != nil {
		return nil, err
	}

	return storedKey, nil
}

type CreateLLMAPIKeyParams struct {
	Name         string
	EncryptedKey []byte
	UserID       string
}
