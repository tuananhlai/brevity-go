package llmapikey

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Crypter interface {
	Encrypt(plainText []byte) []byte
	Decrypt(cipherText []byte) ([]byte, error)
}

// Service defines business logic for managing LLM API keys.
type Service interface {
	ListByUserID(ctx context.Context, userID string) ([]*APIKey, error)
	Create(ctx context.Context, apiKey CreateInput) (*APIKey, error)
}

type serviceImpl struct {
	repo    Repository
	crypter Crypter
}

func NewService(repo Repository, crypter Crypter) Service {
	return &serviceImpl{
		repo:    repo,
		crypter: crypter,
	}
}

func (s *serviceImpl) ListByUserID(ctx context.Context, userID string) ([]*APIKey, error) {
	results, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]*APIKey, 0, len(results))
	for _, result := range results {
		apiKeyBytes, err := s.crypter.Decrypt(result.EncryptedKey)
		if err != nil {
			return nil, err
		}

		apiKey := string(apiKeyBytes)

		res = append(res, &APIKey{
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

func (s *serviceImpl) Create(ctx context.Context, apiKey CreateInput) (*APIKey, error) {
	encryptedKey := s.crypter.Encrypt([]byte(apiKey.Value))

	newAPIKey, err := s.repo.Create(ctx, CreateParams{
		Name:         apiKey.Name,
		EncryptedKey: encryptedKey,
		UserID:       apiKey.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &APIKey{
		ID:            newAPIKey.ID,
		Name:          newAPIKey.Name,
		ValueFirstTen: apiKey.Value[:10],
		ValueLastSix:  apiKey.Value[len(apiKey.Value)-6:],
		UserID:        newAPIKey.UserID,
		CreatedAt:     newAPIKey.CreatedAt,
	}, nil
}

// APIKey is the DTO returned by the service layer.
type APIKey struct {
	ID   uuid.UUID
	Name string
	// ValueFirstTen is the first ten characters of the API key.
	ValueFirstTen string
	// ValueLastSix is the last six characters of the API key.
	ValueLastSix string
	UserID       uuid.UUID
	CreatedAt    time.Time
}

type CreateInput struct {
	Name string
	// Value is the plaintext API key string.
	Value  string
	UserID string
}
