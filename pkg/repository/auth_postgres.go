package repository

import (
	"documentStorage/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (string, error) {
	var login string
	query := fmt.Sprintf("INSERT INTO %s (login, password_hash)"+
		" values ($1, $2) RETURNING login", userTable)
	row := r.db.QueryRow(query, user.Login, user.Password)
	if err := row.Scan(&login); err != nil {
		return "", err
	}

	return login, nil
}
