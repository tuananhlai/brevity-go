package store

import "errors"

var (
	ErrArticleNotFound   = errors.New("article not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
