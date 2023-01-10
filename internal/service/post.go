package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type Post interface {
	CreatePost(post models.Post) error
	GetAllPosts() ([]models.Post, error)
	GetPost(title string) (models.Post, error)
	// DeletePost()
	// UpdatePost()
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

func (s *PostService) GetPost(title string) (models.Post, error) {
	return s.repo.GetPost(title)
}
