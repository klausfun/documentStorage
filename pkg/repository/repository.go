package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type Document interface {
}

type Repository struct {
	Authorization
	Document
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
