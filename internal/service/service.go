package service

import (
	"forum/internal/repository"
)

type (
	Post    interface{}
	Comment interface{}
)

type Service struct {
	Autorization
	Post
	Comment
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Autorization: NewAuthService(repos.Autorization),
	}
}
