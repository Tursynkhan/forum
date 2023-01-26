package service

import (
	"errors"
	"forum/internal/models"
	"forum/internal/repository"
	"strings"
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

var (
	ErrInvalidPost    = errors.New("invalid post")
	ErrPostTitleLen   = errors.New("title length out of range")
	ErrPostContentLen = errors.New("content length out of range")
)

func (s *PostService) CreatePost(post models.Post) (int, error) {
	if isInvalidPost(post) {
		return 0, ErrInvalidPost
	}
	if err := checkPost(post); err != nil {
		return 0, err
	}
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
					posts, err = s.repo.GetPostsByNewest()
					if err != nil {
						return []models.PostInfo{}, nil
					}
				} else if w == "old" {
					posts, err = s.repo.GetPostsByOldest()
					if err != nil {
						return []models.PostInfo{}, nil
					}
				}
			}
		} else if key == "select" {
			for _, w := range val {
				posts, err = s.repo.GetPostByCategory(strings.ReplaceAll(w, "+", " "))
				if err != nil {
					return []models.PostInfo{}, errors.New("post filter service")
				}
			}
		} else {
			return []models.PostInfo{}, err
		}
	}
	return posts, nil
}

func checkPost(post models.Post) error {
	if len(post.Title) > 100 {
		return ErrPostTitleLen
	}

	if len(post.Content) > 1500 {
		return ErrPostContentLen
	}
	return nil
}

func isInvalidPost(post models.Post) bool {
	if strings.ReplaceAll(post.Title, " ", "") == "" {
		return true
	}
	if strings.ReplaceAll(post.Content, " ", "") == "" {
		return true
	}
	return false
}
