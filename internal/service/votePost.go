package service

import (
	"database/sql"
	"errors"
	"forum/internal/models"
	"forum/internal/repository"
)

type VotePostService struct {
	repo repository.VotePost
}

type VotePost interface {
	CreateLikePost(postLike models.PostLike) error
}

func NewVotePostService(repo repository.VotePost) *VotePostService {
	return &VotePostService{repo: repo}
}

func (s *VotePostService) CreateLikePost(postLike models.PostLike) error {
	status, err := s.repo.GetStatusPostLike(postLike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.repo.CreateLikePost(postLike)
		}
		return err
	}
	if status == 1 {
		if err := s.repo.UpdateStatusPostLike(0, postLike); err != nil {
			return err
		}
	} else {
		if err := s.repo.UpdateStatusPostLike(1, postLike); err != nil {
			return err
		}
	}
	return nil
}
