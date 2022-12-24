package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type (
	Autorization interface {
		CreateUser(user models.User) error
	}
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
