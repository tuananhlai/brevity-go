package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
	// qb is a query builder for PostgreSQL
	qb sq.StatementBuilderType
}

var _ Repository = (*Postgres)(nil)

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
