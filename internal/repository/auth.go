package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/tuananhlai/brevity-go/internal/model"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type AuthRepository interface {
	// GetUser returns the user with the given email or username.
	GetUser(ctx context.Context, emailOrUsername string) (*model.AuthUser, error)
	// CreateUser creates a new user and returns the created user.
	CreateUser(ctx context.Context, params CreateUserParams) (*model.AuthUser, error)
	// CreateRefreshToken creates a new refresh token and returns the created refresh token.
	CreateRefreshToken(ctx context.Context, params CreateRefreshTokenParams) (*model.RefreshToken, error)
	// GetRefreshToken returns the information related to the given refresh token.
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
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
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
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
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (r *authRepositoryImpl) CreateRefreshToken(ctx context.Context,
	params CreateRefreshTokenParams,
) (*model.RefreshToken, error) {
	token := &model.RefreshToken{
		Token:     params.Token,
		UserID:    params.UserID,
		ExpiresAt: params.ExpiresAt,
	}

	err := r.db.GetContext(ctx, token,
		`INSERT INTO refresh_tokens (token, user_id, expires_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, token, user_id, expires_at, created_at, revoked_at`,
		params.Token,
		params.UserID,
		params.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *authRepositoryImpl) GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := r.db.GetContext(ctx, &refreshToken,
		`SELECT id, token, user_id, expires_at, created_at, revoked_at
		FROM refresh_tokens WHERE token = $1`, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return &refreshToken, nil
}

type CreateRefreshTokenParams struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

type CreateUserParams struct {
	Email        string
	PasswordHash []byte
	Username     string
}
