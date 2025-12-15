package auth

import "github.com/google/uuid"

// User represents an application user record stored in the database.
type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
}
