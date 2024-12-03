package service

import (
	"documentStorage/models"
	"documentStorage/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (string, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (int, error)
	Logout(token string) error
}

type Document interface {
	Create(meta models.GetDocsResp, fileData []byte, jsonData string) error
	GetListOfDocs(userId int, docInput models.GetDocsInput) ([]models.GetDocsResp, error)
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
