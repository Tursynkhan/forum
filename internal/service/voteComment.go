package service

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
)

type VoteCommentService struct {
	repo repository.VoteComment
}

type VoteComment interface {
	CreateLikeComment(comment models.CommentLike) error
	CreateDisLikeComment(comment models.CommentLike) error
	GetCommentLikesByCommentID(id int) (int, error)
	GetCommentDislikesByCommentID(id int) (int, error)
}

func NewVoteCommentService(repo repository.VoteComment) *VoteCommentService {
	return &VoteCommentService{repo: repo}
}

func (s *VoteCommentService) CreateLikeComment(comment models.CommentLike) error {
	status, err := s.repo.GetStatusCommentLike(comment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.repo.CreateLikeComment(comment)
		}
	}
	if status == 1 {
		if err := s.repo.UpdateStatusCommentLike(0, comment); err != nil {
			return fmt.Errorf("service : UpdateStatusCommentLike : %w", err)
		}
	} else {
		if err := s.repo.UpdateStatusCommentLike(1, comment); err != nil {
			return fmt.Errorf("service : UpdateStatusCommentLike : %w", err)
		}
	}
	return nil
}

func (s *VoteCommentService) CreateDisLikeComment(comment models.CommentLike) error {
	status, err := s.repo.GetStatusCommentLike(comment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.repo.CreateDisLikeComment(comment)
		}
	}
	if status == 1 {
		if err := s.repo.UpdateStatusCommentLike(-1, comment); err != nil {
			return fmt.Errorf("service : UpdateStatusCommentLike : %w", err)
		}
	} else if status == -1 {
		if err := s.repo.UpdateStatusCommentLike(0, comment); err != nil {
			return fmt.Errorf("service : UpdateStatusCommentLike : %w", err)
		}
	} else {
		if err := s.repo.UpdateStatusCommentLike(-1, comment); err != nil {
			return fmt.Errorf("service : UpdateStatusCommentLike : %w", err)
		}
	}
	return nil
}

func (s *VoteCommentService) GetCommentLikesByCommentID(id int) (int, error) {
	likes, err := s.repo.GetCommentLikesByCommentID(id)
	if err != nil {
		return 0, fmt.Errorf("service : GetCommentLikesByCommentID : %w", err)
	}
	return likes, nil
}

func (s *VoteCommentService) GetCommentDislikesByCommentID(id int) (int, error) {
	dislikes, err := s.repo.GetCommentDislikesByCommentID(id)
	if err != nil {
		return 0, fmt.Errorf("service : GetCommentDislikesByCommentID : %w", err)
	}
	return dislikes, nil
}
