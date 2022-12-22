package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type (
	Autorization interface {
		CreateUser(models.User) (int, error)
		GenerateToken(username, password) (string, error)
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
	return &Service{}
}
