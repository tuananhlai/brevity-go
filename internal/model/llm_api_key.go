package model

import (
	"time"

	"github.com/google/uuid"
)

type LLMAPIKey struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	EncryptedKey []byte    `db:"encrypted_key"`
	UserID       uuid.UUID `db:"user_id"`
	CreatedAt    time.Time `db:"created_at"`
}
