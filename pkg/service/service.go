package service

import (
	"documentStorage/models"
	"documentStorage/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (string, error)
	GenerateToken(login, password string) (string, error)
}

type Document interface {
}

type Service struct {
	Authorization
	Document
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos),
	}
}
