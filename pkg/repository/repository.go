package repository

import (
	"documentStorage/models"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Authorization interface {
	CreateUser(user models.User) (string, error)
	GetUser(login, password string) (models.User, error)
	CreateToken(token string) error
	GetToken(token string) (string, error)
}

type Document interface {
	Create(userId int, meta models.Document, fileData []byte, jsonData string) error
}

type Repository struct {
	Authorization
	Document
}

func NewRepository(db *sqlx.DB, redis *redis.Client) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db, redis),
		Document:      NewDocumentPostgres(db, redis),
	}
}
