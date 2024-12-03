package service

import (
	"documentStorage/models"
	"documentStorage/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user models.User) (string, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (int, error)
	Logout(token string) error
}

type Document interface {
	Create(meta models.GetDocsResp, fileData []byte, jsonData string) error
	GetList(userId int, docInput models.GetDocsInput) ([]models.GetDocsResp, error)
	GetById(docId int) (models.GetDoc, error)
	Delete(docId int) error
}

type Service struct {
	Authorization
	Document
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos),
		Document:      NewDocumentService(repos),
	}
}
