package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tuananhlai/brevity-go/internal/repository"
)

type LLMAPIKeyService interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*LLMAPIKey, error)
	Create(ctx context.Context, apiKey LLMAPIKeyCreateParams) (*LLMAPIKey, error)
}

type llmAPIKeyServiceImpl struct {
	repo              repository.LLMAPIKeyRepository
	encryptionService EncryptionService
}

func NewLLMAPIKeyService(repo repository.LLMAPIKeyRepository, encryptionService EncryptionService) LLMAPIKeyService {
	return &llmAPIKeyServiceImpl{
		repo:              repo,
		encryptionService: encryptionService,
	}
}

func (s *llmAPIKeyServiceImpl) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*LLMAPIKey, error) {
	results, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]*LLMAPIKey, 0, len(results))
	for _, result := range results {
		apiKeyBytes, err := s.encryptionService.Decrypt(result.EncryptedKey)
		if err != nil {
			return nil, err
		}

		apiKey := string(apiKeyBytes)

		res = append(res, &LLMAPIKey{
			ID:            result.ID,
			Name:          result.Name,
			ValueFirstTen: apiKey[:10],
			ValueLastSix:  apiKey[len(apiKey)-6:],
			UserID:        result.UserID,
			CreatedAt:     result.CreatedAt,
		})
	}

	return res, nil
}

func (s *llmAPIKeyServiceImpl) Create(ctx context.Context, apiKey LLMAPIKeyCreateParams) (*LLMAPIKey, error) {
	encryptedKey, err := s.encryptionService.Encrypt([]byte(apiKey.Value))
	if err != nil {
		return nil, err
	}

	newAPIKey, err := s.repo.Create(ctx, repository.LLMAPIKeyCreateParams{
		Name:         apiKey.Name,
		EncryptedKey: encryptedKey,
		UserID:       apiKey.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &LLMAPIKey{
		ID:            newAPIKey.ID,
		Name:          newAPIKey.Name,
		ValueFirstTen: apiKey.Value[:10],
		ValueLastSix:  apiKey.Value[len(apiKey.Value)-6:],
		UserID:        newAPIKey.UserID,
		CreatedAt:     newAPIKey.CreatedAt,
	}, nil
}

type LLMAPIKey struct {
	ID   uuid.UUID
	Name string
	// ValueFirstTen is the first ten characters of the API key.
	ValueFirstTen string
	// ValueLastSix is the last six characters of the API key.
	ValueLastSix string
	UserID       uuid.UUID
	CreatedAt    time.Time
}

type LLMAPIKeyCreateParams struct {
	Name string
	// Value is the plaintext API key string.
	Value  string
	UserID uuid.UUID
}
