package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ContentFormat string

const (
	ContentFormatHTML ContentFormat = "html"
)

type Article struct {
	ID               uuid.UUID     `db:"id"`
	Slug             string        `db:"slug"`
	Title            string        `db:"title"`
	Description      string        `db:"description"`
	PlaintextContent string        `db:"plaintext_content"`
	Content          string        `db:"content"`
	ContentFormat    ContentFormat `db:"content_format"`
	AuthorID         uuid.UUID     `db:"author_id"`
	CreatedAt        time.Time     `db:"created_at"`
	UpdatedAt        time.Time     `db:"updated_at"`
}

type ArticlePreview struct {
	ID                uuid.UUID      `db:"id"`
	Slug              string         `db:"slug"`
	Title             string         `db:"title"`
	Description       string         `db:"description"`
	AuthorID          uuid.UUID      `db:"author_id"`
	AuthorUsername    string         `db:"author_username"`
	AuthorDisplayName sql.NullString `db:"author_display_name"`
	AuthorAvatarURL   sql.NullString `db:"author_avatar_url"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
}

type ArticleDetails struct {
	ID                uuid.UUID      `db:"id"`
	Slug              string         `db:"slug"`
	Title             string         `db:"title"`
	Content           string         `db:"content"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
	AuthorID          uuid.UUID      `db:"author_id"`
	AuthorUsername    string         `db:"author_username"`
	AuthorDisplayName sql.NullString `db:"author_display_name"`
	AuthorAvatarURL   sql.NullString `db:"author_avatar_url"`
}

type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
}

// StoredAPIKey represents the persisted API key with encrypted value.
type StoredAPIKey struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	EncryptedKey []byte    `db:"encrypted_key"`
	UserID       uuid.UUID `db:"user_id"`
	CreatedAt    time.Time `db:"created_at"`
}
