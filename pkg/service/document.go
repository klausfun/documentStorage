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

func (s *DocumentService) GetList(userId int, docInput models.GetDocsInput) ([]models.GetDocsResp, error) {
	return s.repo.GetList(userId, docInput)
}

func (s *DocumentService) GetById(docId int) (models.GetDoc, error) {
	return s.repo.GetById(docId)
}

func (s *DocumentService) Delete(docId int) error {
	return s.repo.Delete(docId)
}
