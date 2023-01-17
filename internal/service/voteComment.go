package service

import (
	"database/sql"
	"errors"
	"forum/internal/models"
	"forum/internal/repository"
)

type VoteCommentService struct {
	repo repository.VoteComment
}

type VoteComment interface {
	CreateLikeComment(comment models.CommentLike) error
	CreateDisLikeComment(comment models.CommentLike) error
	GetAllLikesByCommentId(commentId int) (int, error)
	GetAllDislikesByCommentId(commentId int) (int, error)
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
		if err := s.repo.UpdateStatusCommentLike(-1, comment); err != nil {
			return err
		}
	} else if status == -1 {
		if err := s.repo.UpdateStatusCommentLike(0, comment); err != nil {
			return err
		}
	} else {
		if err := s.repo.UpdateStatusCommentLike(-1, comment); err != nil {
			return err
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
			return err
		}
	} else if status == -1 {
		if err := s.repo.UpdateStatusCommentLike(0, comment); err != nil {
			return err
		}
	} else {
		if err := s.repo.UpdateStatusCommentLike(-1, comment); err != nil {
			return err
		}
	}
	return nil
}

func (s *VoteCommentService) GetAllLikesByCommentId(commentId int) (int, error) {
	return s.repo.GetAllLikesByCommentId(commentId)
}

func (s *VoteCommentService) GetAllDislikesByCommentId(commentId int) (int, error) {
	return s.repo.GetAllDislikesByCommentId(commentId)
}
