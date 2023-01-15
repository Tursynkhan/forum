package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type Comment interface {
	CreateComment(comment models.Comment) error
	GetAllComments(postId int) ([]models.Comment, error)
}

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(comment models.Comment) error {
	return s.repo.CreateComment(comment)
}

func (s *CommentService) GetAllComments(postId int) ([]models.Comment, error) {
	return s.repo.GetAllComments(postId)
}
