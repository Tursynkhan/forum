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
	GetPostByFilter(query map[string][]string) ([]models.PostInfo, error)
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

func (s *PostService) GetPostByFilter(query map[string][]string) ([]models.PostInfo, error) {
	var posts []models.PostInfo
	var err error
	for key, val := range query {
		if key == "like" {
			for _, w := range val {
				if w == "most" {
					posts, err = s.repo.GetPostsByMostLikes()
					if err != nil {
						return []models.PostInfo{}, err
					}
				} else if w == "least" {
					posts, err = s.repo.GetPostsByLeastLikes()
					if err != nil {
						return []models.PostInfo{}, err
					}
				}
			}
		} else if key == "time" {
			for _, w := range val {
				if w == "new" {
				} else if w == "old" {
				}
			}
		} else if key == "select" {
			for _, w := range val {
				posts, err = s.repo.GetPostByCategory(w)
				if err != nil {
					return []models.PostInfo{}, err
				}
			}
		}
	}
	return posts, nil
}
