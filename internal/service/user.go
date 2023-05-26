package service

import (
	"database/sql"
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
	ApprovedPost(postId int) error
	DeclinePost(postId int) error
	UpdateReportOfPost(postId int, report string) error
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
				} else if w == "request" {
					posts, err = s.repo.GetPostsToApprove()
					if err != nil {
						return []models.PostInfo{}, err
					}
				} else if w == "irrelevant" {
					posts, err = s.repo.GetPostByReportName("irrelevant")
					if err != nil {
						return []models.PostInfo{}, err
					}
				} else if w == "obscene" {
					posts, err = s.repo.GetPostByReportName("obscene")
					if err != nil {
						return []models.PostInfo{}, err
					}
				} else if w == "illegal" {
					posts, err = s.repo.GetPostByReportName("illegal")
					if err != nil {
						return []models.PostInfo{}, err
					}
				} else if w == "insulting" {
					posts, err = s.repo.GetPostByReportName("insulting")
					if err != nil {
						return []models.PostInfo{}, err
					}
				}
			}
		} else {
			return []models.PostInfo{}, err
		}
	}

	return posts, nil
}

func (s *UserService) GetProfileByUsername(username string) (models.ProfileUser, error) {
	profile, err := s.repo.GetProfileByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ProfileUser{}, fmt.Errorf("service : GetProfileByUsername :%w", err)
		}
		return models.ProfileUser{}, fmt.Errorf("service : GetProfileByUsername :%w", err)
	}
	return profile, nil
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
		return fmt.Errorf("service: DeleteCategoryByName : %w", err)
	}
	return s.repo.DeleteCategoryByName(name)
}
func (s *UserService) UpdateUserRole(username string, roleId int) error {

	return s.repo.UpdateUserRole(username, roleId)
}

func (s *UserService) ApprovedPost(postId int) error {
	return s.repo.ApprovedPost(postId)
}
func (s *UserService) DeclinePost(postId int) error {
	return s.repo.DeclinePost(postId)
}
func (s *UserService) UpdateReportOfPost(postId int,report string) error {
	return s.repo.UpdateReportOfPost(postId,report)
}
