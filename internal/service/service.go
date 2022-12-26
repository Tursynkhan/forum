package service

import (
	"time"

	"forum/internal/models"
	"forum/internal/repository"
)

type (
	Autorization interface {
		CreateUser(user models.User) error
		GenerateToken(username, password string) (string, time.Time, error)
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
