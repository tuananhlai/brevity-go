package model

import (
	"time"

	"github.com/google/uuid"
)

type DigitalAuthor struct {
	ID           uuid.UUID `db:"id"`
	OwnerID      uuid.UUID `db:"owner_id"`
	DisplayName  string    `db:"display_name"`
	SystemPrompt string    `db:"system_prompt"`
	AvatarURL    string    `db:"avatar_url"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
