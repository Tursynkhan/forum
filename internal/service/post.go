package service

import (
	"errors"
	"fmt"
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
	GetPostByFilter(query map[string][]string, user models.User) ([]models.PostInfo, error)
	GetLenAllPost() (int, error)
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
	posts, err := s.repo.GetAllPosts()
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("service : GetAllPosts: %w", err)
	}
	return posts, nil
}

func (s *PostService) GetPost(id int) (models.PostInfo, error) {
	post, err := s.repo.GetPost(id)
	if err != nil {
		return models.PostInfo{}, fmt.Errorf("service : GetPost: %w", err)
	}
	return post, nil
}

func (s *PostService) CreatePostCategory(id int, categories []string) error {
	if err := s.repo.CreatePostCategory(id, categories); err != nil {
		return fmt.Errorf("service : CreatePostCategory %w:", err)
	}
	return nil
}

func (s *PostService) GetAllCategories() ([]models.Category, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		return []models.Category{}, fmt.Errorf("service : GetAllCategories : %w", err)
	}
	return categories, nil
}

func (s *PostService) GetPostByFilter(query map[string][]string, user models.User) ([]models.PostInfo, error) {
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
		} else if key == "my" {
			for _, w := range val {
				if w == "post" {
					posts, err = s.repo.GetMyPosts(user)
					if err != nil {
						return []models.PostInfo{}, nil
					}
				} else if w == "like" {
					posts, err = s.repo.GetMyLikedPosts(user)
					if err != nil {
						return []models.PostInfo{}, nil
					}
				}
			}
		} else if key == "select" {
			for _, w := range val {
				posts, err = s.repo.GetPostByCategory(w)
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

func (s *PostService) GetLenAllPost() (int, error) {
	count, err := s.repo.GetLenAllPost()
	if err != nil {
		return 0, fmt.Errorf("service : post : GetLenAllPost %w", err)
	}
	return count, nil
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
