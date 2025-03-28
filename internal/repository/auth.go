package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/tuananhlai/brevity-go/internal/model"
)

type AuthRepository interface {
	// GetUser returns the user with the given email or username.
	GetUser(ctx context.Context, emailOrUsername string) (*model.AuthUser, error)
	// CreateUser creates a new user and returns the created user.
	CreateUser(ctx context.Context, params CreateUserParams) (*model.AuthUser, error)
}

type authRepositoryImpl struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}

func (r *authRepositoryImpl) GetUser(ctx context.Context, emailOrUsername string) (*model.AuthUser, error) {
	var user model.AuthUser
	err := r.db.GetContext(ctx, &user,
		`SELECT id, username, email, password_hash
		FROM users 
		WHERE (email = $1 OR username = $2) LIMIT 1`,
		emailOrUsername,
		emailOrUsername)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepositoryImpl) CreateUser(ctx context.Context, params CreateUserParams) (*model.AuthUser, error) {
	user := &model.AuthUser{
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
		return nil, err
	}

	return user, nil
}

type CreateUserParams struct {
	Email        string
	PasswordHash []byte
	Username     string
}
