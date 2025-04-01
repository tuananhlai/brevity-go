package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"

	"github.com/tuananhlai/brevity-go/internal/model"
	"github.com/tuananhlai/brevity-go/internal/otelsdk"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type AuthRepository interface {
	// GetUser returns the user with the given email or username.
	GetUser(ctx context.Context, emailOrUsername string) (*model.AuthUser, error)
	// CreateUser creates a new user and returns the created user.
	CreateUser(ctx context.Context, params CreateUserParams) (*model.AuthUser, error)
}

type authRepositoryImpl struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepositoryImpl{db: db, tracer: otelsdk.Tracer("repository.AuthRepository")}
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

type CreateUserParams struct {
	Email        string
	PasswordHash []byte
	Username     string
}
