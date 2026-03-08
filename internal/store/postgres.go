package store

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
	// qb is a query builder for PostgreSQL
	qb sq.StatementBuilderType
}

func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
