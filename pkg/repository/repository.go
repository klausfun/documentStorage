package repository

type Authorization interface {
}

type Document interface {
}

type Repository struct {
	Authorization
	Document
}

func NewRepository() *Repository {
	return &Repository{}
}
