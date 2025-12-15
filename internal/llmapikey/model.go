package llmapikey

import (
	"time"

	"github.com/google/uuid"
)

// StoredAPIKey represents the persisted API key with encrypted value.
type StoredAPIKey struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	EncryptedKey []byte    `db:"encrypted_key"`
	UserID       uuid.UUID `db:"user_id"`
	CreatedAt    time.Time `db:"created_at"`
}
