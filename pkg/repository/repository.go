package repository

import (
	"documentStorage/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (string, error)
	GetUser(login, password string) (models.User, error)
}

type Document interface {
}

type Repository struct {
	Authorization
	Document
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
