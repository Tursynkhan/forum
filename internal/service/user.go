package service

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
)

var ErrInvalidQuery = errors.New("Invalid query request")

type User interface {
	GetPostByUsername(username string, query map[string][]string) ([]models.PostInfo, error)
	GetProfileByUsername(username string) (models.ProfileUser, error)
	GetAllRoles() ([]models.Role, error)
	CreateCategory(category string) error
	DeleteCategoryById(categoryId int) error
	UpdateUserRole(username string, roleId int) error
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
	// search, ok := query["posts"]
	// if !ok {
	// 	return nil, fmt.Errorf("service: GetPostByUserName: %w", ErrInvalidQuery)
	// }
	for key, val := range query {
		if key == "posts" {
			for _, w := range val {
				if w == "created" {
					posts, err = s.repo.GetPostByUsername(username)
					if err != nil {
						return []models.PostInfo{}, err
					}
				} else if w == "liked" {
					posts, err = s.repo.GetLikedPostByUsername(username)
					if err != nil {
						return []models.PostInfo{}, err
					}
				} else if w == "commented" {
					posts, err = s.repo.GetCommentedPostByUsername(username)
					if err != nil {
						return []models.PostInfo{}, err
					}
				}
			}
		} else {
			return []models.PostInfo{}, err
		}
	}
	// switch strings.Join(search, "") {
	// case "created":
	// 	posts, err = s.repo.GetPostByUsername(username)
	// case "liked":
	// 	posts, err = s.repo.GetLikedPostByUsername(username)
	// case "commented":
	// 	posts, err = s.repo.GetCommentedPostByUsername(username)
	// default:
	// 	return nil, fmt.Errorf("service: GetPostByUsernameL %w", err)
	// }
	return posts, nil
}

func (s *UserService) GetProfileByUsername(username string) (models.ProfileUser, error) {
	return s.repo.GetProfileByUsername(username)
}

func (s *UserService) GetAllRoles() ([]models.Role, error) {
	return s.repo.GetAllRoles()
}
func (s *UserService) CreateCategory(category string) error {
	return s.repo.CreateCategory(category)
}
func (s *UserService) DeleteCategoryById(categoryId int) error {
	name, err := s.repo.GetNameCategoryById(categoryId)
	if err != nil {
		fmt.Errorf("service: DeleteCategoryByName : %w", err)
	}
	return s.repo.DeleteCategoryByName(name)
}
func (s *UserService) UpdateUserRole(username string, roleId int) error {

	return s.repo.UpdateUserRole(username, roleId)
}
