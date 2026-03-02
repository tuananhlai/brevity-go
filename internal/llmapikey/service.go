package llmapikey

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tuananhlai/brevity-go/internal/store"
)

// Crypter defines the interface for encrypting and decrypting data.
type Crypter interface {
	Encrypt(plainText []byte) []byte
	Decrypt(cipherText []byte) ([]byte, error)
}

// Manager handles LLM API key operations, encapsulating encryption/decryption and DTO mapping logic.
type Manager struct {
	store   store.Store
	crypter Crypter
}

// NewManager creates a new LLM API key manager.
func NewManager(store store.Store, crypter Crypter) *Manager {
	return &Manager{
		store:   store,
		crypter: crypter,
	}
}

// ListByUserID returns all API keys for the given user with masked values.
func (m *Manager) ListByUserID(ctx context.Context, userID string) ([]*APIKey, error) {
	results, err := m.store.ListLLMAPIKeysByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]*APIKey, 0, len(results))
	for _, result := range results {
		apiKeyBytes, err := m.crypter.Decrypt(result.EncryptedKey)
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

// Create encrypts and stores a new API key, returning the created key with masked values.
func (m *Manager) Create(ctx context.Context, apiKey CreateInput) (*APIKey, error) {
	encryptedKey := m.crypter.Encrypt([]byte(apiKey.Value))

	newAPIKey, err := m.store.CreateLLMAPIKey(ctx, store.CreateLLMAPIKeyParams{
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

// APIKey is the DTO returned by the manager.
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

// CreateInput is the input for creating a new API key.
type CreateInput struct {
	Name string
	// Value is the plaintext API key string.
	Value  string
	UserID string
}
