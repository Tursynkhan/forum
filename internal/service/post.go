package service

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"io"
	"os"
	"strconv"
	"strings"
)

type Post interface {
	CreatePost(post models.Post) (int, error)
	GetAllPosts() ([]models.PostInfo, error)
	GetPost(id int) (models.PostInfo, error)
	SaveImageForPost(post models.Post) error
	CreatePostCategory(id int, categories []string) error
	EditPostCategory(id int, categories []string) error
	GetAllCategories() ([]models.Category, error)
	GetPostByFilter(query map[string][]string, user models.User) ([]models.PostInfo, error)
	DeletePost(post models.PostInfo, user models.User) error
	EditPost(oldPost models.PostInfo, newPost models.Post, user models.User) error
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
	ErrInvalidType    = errors.New("the provided file format is not allowed")
	ErrInvalidUser    = errors.New("invalid user")
)

func (s *PostService) CreatePost(post models.Post) (int, error) {
	if isInvalidPost(post) {
		return 0, ErrInvalidPost
	}
	if err := checkPost(post); err != nil {
		return 0, err
	}

	id, err := s.repo.CreatePost(post)
	if err != nil {
		return 0, fmt.Errorf("service: createPost: %w", err)
	}
	post.ID = id
	if err := s.SaveImageForPost(post); err != nil {
		if errors.Is(err, ErrInvalidType) {
			return 0, err
		}
		if err := s.repo.DeletePostById(post.ID); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("service: SaveImageForPost: %w", err)
	}
	return id, nil
}

func (s *PostService) DeletePost(post models.PostInfo, user models.User) error {

	if user.Username != post.Author && user.RoleID != 3 && user.RoleID != 4 {
		print(user.Username, post.Author)
		return fmt.Errorf("service : DeletePost : %w: can't delete post", ErrInvalidUser)
	}

	if err := s.repo.DeletePostById(post.ID); err != nil {
		return fmt.Errorf("repo : DeletePost: %w", err)
	}
	return nil
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

func (s *PostService) EditPostCategory(id int, categories []string) error {
	if err := s.repo.EditPostCategory(id, categories); err != nil {
		return fmt.Errorf("service : EditPostCategory %w:", err)
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
		} else if key == "tag" {
			for _, w := range val {
				id, err := strconv.Atoi(w)
				if err != nil {
					return []models.PostInfo{}, errors.New("service: postfilters: can't convert")
				}
				posts, err = s.repo.GetPostByCategory(id)
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

func (s *PostService) SaveImageForPost(post models.Post) error {
	if err := os.MkdirAll(fmt.Sprintf("./ui/static/upload/%d", post.ID), os.ModePerm); err != nil {
		return err
	}
	for _, fileHeader := range post.Files {
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		imageSplit := strings.Split(fileHeader.Filename, ".")
		if !validImageType(imageSplit[len(imageSplit)-1]) {
			return ErrInvalidType
		}

		f, err := os.Create(fmt.Sprintf("./ui/static/upload/%d/%s", post.ID, fileHeader.Filename))
		if err != nil {
			return err
		}

		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
		if err := s.repo.SaveImageForPost(post.ID, fmt.Sprintf("/static/upload/%d/%s", post.ID, fileHeader.Filename)); err != nil {
			return fmt.Errorf("service : SaveImageForPost : %w", err)
		}
	}
	return nil
}

func (s *PostService) EditPost(oldPost models.PostInfo, newPost models.Post, user models.User) error {
	if user.Username != oldPost.Author {
		return fmt.Errorf("service : EditPost : %w : you can't edit post", ErrInvalidUser)
	}
	if isInvalidPost(newPost) {
		return ErrInvalidPost
	}
	if err := checkPost(newPost); err != nil {
		return err
	}
	if err := s.repo.EditPost(newPost, oldPost.ID); err != nil {
		return fmt.Errorf("service : EditPost : %w", err)
	}
	return nil
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

func validImageType(imageType string) bool {
	validImageType := []string{"jpeg", "jpg", "png", "gif"}
	for _, t := range validImageType {
		if t == imageType {
			return true
		}
	}
	return false
}
