package store

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type PostgresStore struct {
	db *sqlx.DB
	// qb is a query builder for PostgreSQL
	qb sq.StatementBuilderType
}

var _ Store = (*PostgresStore)(nil)

func NewPostgresStore(db *sqlx.DB) *PostgresStore {
	return &PostgresStore{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
