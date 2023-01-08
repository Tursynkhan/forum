package service

import (
	"forum/internal/repository"
)

type (
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
		Post:         NewPostService(repos.Post),
	}
}
