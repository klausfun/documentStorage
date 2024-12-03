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

func (s *DocumentService) Create(meta models.GetDocsResp,
	fileData []byte, jsonData string) error {

	return s.repo.Create(meta, fileData, jsonData)
}

func (s *DocumentService) GetListOfDocs(userId int, docInput models.GetDocsInput) ([]models.GetDocsResp, error) {
	return s.repo.GetListOfDocs(userId, docInput)
}
