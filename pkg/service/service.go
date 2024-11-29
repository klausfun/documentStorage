package service

import "documentStorage/pkg/repository"

type Authorization interface {
}

type Document interface {
}

type Service struct {
	Authorization
	Document
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
