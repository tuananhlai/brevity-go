package repository

import "github.com/jmoiron/sqlx"

type Postgres struct {
	db *sqlx.DB
}

var _ Repository = (*Postgres)(nil)

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}
