package service

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
)

type Comment interface {
	CreateComment(comment models.Comment) error
	GetAllComments(postId int) ([]models.Comment, error)
	GetCommentById(id int) (models.Comment, error)
}

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(comment models.Comment) error {
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
