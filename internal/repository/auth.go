package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/tuananhlai/brevity-go/internal/model"
)

type AuthRepository interface {
	// GetUser returns the user with the given email and password hash.
	GetUser(ctx context.Context, email string, passwordHash string) (*model.AuthUser, error)
	// CreateUser creates a new user and returns the created user.
	CreateUser(ctx context.Context, email string, passwordHash string) (*model.AuthUser, error)
}

type authRepositoryImpl struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}

func (r *authRepositoryImpl) GetUser(ctx context.Context, email string, passwordHash string) (*model.AuthUser, error) {
	var user *model.AuthUser
	err := r.db.GetContext(ctx, user, "SELECT id, username, email FROM users WHERE email = $1 AND password_hash = $2", email, passwordHash)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *authRepositoryImpl) CreateUser(ctx context.Context, email string, passwordHash string) (*model.AuthUser, error) {
	user := &model.AuthUser{
		Email: email,
		PasswordHash: passwordHash,
	}

	err := r.db.GetContext(ctx, user, "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, username, email", email, passwordHash)
	if err != nil {
		return nil, err
	}

	return user, nil
}
