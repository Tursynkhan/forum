package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type Post interface {
	CreatePost(post models.Post) error
	GetAllPosts() ([]models.Post, error)
	GetPost(id int) (models.Post, error)
}

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(post models.Post) error {
	return s.repo.CreatePost(post)
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	return s.repo.GetAllPosts()
}

func (s *PostService) GetPost(id int) (models.Post, error) {
	return s.repo.GetPost(id)
}
