package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type Post interface {
	CreatePost(post models.Post) (int, error)
	GetAllPosts() ([]models.PostInfo, error)
	GetPost(id int) (models.PostInfo, error)
	CreatePostCategory(id int, categories []string) error
	GetAllCategories() ([]models.Category, error)
}

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(post models.Post) (int, error) {
	return s.repo.CreatePost(post)
}

func (s *PostService) GetAllPosts() ([]models.PostInfo, error) {
	return s.repo.GetAllPosts()
}

func (s *PostService) GetPost(id int) (models.PostInfo, error) {
	return s.repo.GetPost(id)
}

func (s *PostService) CreatePostCategory(id int, categories []string) error {
	return s.repo.CreatePostCategory(id, categories)
}

func (s *PostService) GetAllCategories() ([]models.Category, error) {
	return s.repo.GetAllCategories()
}
