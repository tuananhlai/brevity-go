package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// Repository defines the read/write operations for users.
type Repository interface {
	// GetUser returns the user with the given email or username.
	GetUser(ctx context.Context, emailOrUsername string) (*User, error)
	// CreateUser creates a new user and returns the created user.
	CreateUser(ctx context.Context, params CreateUserParams) (*User, error)
	// GetUserByID returns the user with the given ID.
	GetUserByID(ctx context.Context, userID string) (*User, error)
}

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) GetUser(ctx context.Context, emailOrUsername string) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user,
		`SELECT id, username, email, password_hash
		FROM users 
		WHERE (email = $1 OR username = $2) LIMIT 1`,
		emailOrUsername,
		emailOrUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *repositoryImpl) CreateUser(ctx context.Context, params CreateUserParams) (*User, error) {
	user := &User{
		Email:    params.Email,
		Username: params.Username,
	}

	err := r.db.GetContext(ctx, user,
		`INSERT INTO users (email, username, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING id, username, email`,
		params.Email,
		params.Username,
		params.PasswordHash,
	)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (r *repositoryImpl) GetUserByID(ctx context.Context, userID string) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user,
		`SELECT id, username, email, password_hash
		FROM users 
		WHERE id = $1`,
		userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

type CreateUserParams struct {
	Email        string
	PasswordHash []byte
	Username     string
}
