package service

import "forum/internal/repository"

type (
	Autorization interface{}
	Post         interface{}
	Comment      interface{}
)

type Service struct {
	Autorization
	Post
	Comment
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
