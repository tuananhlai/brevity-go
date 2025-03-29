package model

import "github.com/google/uuid"

type AuthUser struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
}
