package repository

import (
	"context"
	"database/sql"
	"errors"
)

func (r *Postgres) GetUser(ctx context.Context, emailOrUsername string) (*User, error) {
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

func (r *Postgres) CreateUser(ctx context.Context, params CreateUserParams) (*User, error) {
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

func (r *Postgres) GetUserByID(ctx context.Context, userID string) (*User, error) {
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
