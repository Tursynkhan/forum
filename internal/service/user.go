package service

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"strings"
)

type User interface {
	GetPostByUsername(username string, query map[string][]string) ([]models.PostInfo, error)
	GetUserByUsername(username string) (models.User, error)
}

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetPostByUsername(username string, query map[string][]string) ([]models.PostInfo, error) {
	var (
		posts []models.PostInfo
		err   error
	)
	search, ok := query["posts"]
	if !ok {
		return nil, fmt.Errorf("Invalid query request")
	}
	switch strings.Join(search, "") {
	case "created":
		posts, err = s.repo.GetPostByUsername(username)
	case "liked":
		posts, err = s.repo.GetLikedPostByUsername(username)
	case "commented":
		posts, err = s.repo.GetCommentPostByUsername(username)
	default:
		return nil, fmt.Errorf("service: GetPostByUsernameL %w", err)
	}
	return posts, nil
}

func (s *UserService) GetUserByUsername(username string) (models.User, error) {
	return s.repo.GetUserByUsername(username)
}
