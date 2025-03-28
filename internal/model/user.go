package model

type AuthUser struct {
	ID string `db:"id"`
	Username string `db:"username"`
	Email string `db:"email"`
}
