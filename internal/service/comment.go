package service

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"strings"
)

var ErrPostNotexist = errors.New("Post does not exist")
var ErrInvalidComment = errors.New("invalid comment")

type Comment interface {
	CreateComment(comment models.Comment) error
	GetAllComments(postId int) ([]models.Comment, error)
	GetCommentById(id int) (models.Comment, error)
	DeleteComment(comment models.Comment) error
	EditComment(comment models.Comment) error
}

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(comment models.Comment) error {
	_, err := s.repo.GetPostById(comment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPostNotexist
		}
		return fmt.Errorf("service : comment : CreateComment : %w", err)
	}
	if err := s.repo.CreateComment(comment); err != nil {
		return fmt.Errorf("service : comment : CreateComment : %w", err)
	}
	return nil
}

func (s *CommentService) GetAllComments(postId int) ([]models.Comment, error) {
	comments, err := s.repo.GetAllComments(postId)
	if err != nil {
		return []models.Comment{}, fmt.Errorf("service : comment : GetAllComments : %w", err)
	}
	return comments, nil
}

func (s *CommentService) GetCommentById(id int) (models.Comment, error) {
	comment, err := s.repo.GetCommentById(id)
	if err != nil {
		return models.Comment{}, fmt.Errorf("service : comment : GetCommentById : %w", err)
	}
	return comment, nil
}
func (s *CommentService) DeleteComment(comment models.Comment) error {
	return s.repo.DeleteComment(comment)
}

func (s *CommentService) EditComment(comment models.Comment) error {
	if strings.ReplaceAll(comment.Content, " ", "") == "" {
		return fmt.Errorf("service: edit comment: %w", ErrInvalidComment)
	}
	return s.repo.EditComment(comment)
}
