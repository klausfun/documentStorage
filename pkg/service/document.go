package service

import (
	"documentStorage/models"
	"documentStorage/pkg/repository"
)

type DocumentService struct {
	repo repository.Document
}

func NewDocumentService(repo repository.Document) *DocumentService {
	return &DocumentService{
		repo: repo,
	}
}

func (s *DocumentService) Create(userId int, meta models.Document,
	fileData []byte, jsonData string) error {

	return s.repo.Create(userId, meta, fileData, jsonData)
}
