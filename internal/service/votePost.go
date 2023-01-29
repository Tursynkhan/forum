package service

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
)

type VotePostService struct {
	repo repository.VotePost
}

type VotePost interface {
	CreateLikePost(postLike models.PostLike) error
	CreateDisLikePost(postLike models.PostLike) error
	GetAllLikesByPostId(postId int) (int, error)
	GetAllDislikesByPostId(postId int) (int, error)
}

func NewVotePostService(repo repository.VotePost) *VotePostService {
	return &VotePostService{repo: repo}
}

func (s *VotePostService) CreateLikePost(postLike models.PostLike) error {
	_, err := s.repo.GetPostById(postLike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPostNotexist
		}
		return fmt.Errorf("service: CreateLikePost: %w", err)
	}
	status, err := s.repo.GetStatusPostLike(postLike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.repo.CreateLikePost(postLike)
		}
		return err
	}
	if status == 1 {
		if err := s.repo.UpdateStatusPostLike(0, postLike); err != nil {
			return fmt.Errorf("service : post : UpdateStatusPostLike : %w", err)
		}
	} else {
		if err := s.repo.UpdateStatusPostLike(1, postLike); err != nil {
			return fmt.Errorf("service : post : UpdateStatusPostLike : %w", err)
		}
	}
	return nil
}

func (s *VotePostService) CreateDisLikePost(postLike models.PostLike) error {
	_, err := s.repo.GetPostById(postLike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPostNotexist
		}
		return fmt.Errorf("service : CreateDislikePost: %w", err)
	}
	status, err := s.repo.GetStatusPostLike(postLike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.repo.CreateDisLikePost(postLike)
		}
	}
	if status == 1 {
		if err := s.repo.UpdateStatusPostLike(-1, postLike); err != nil {
			return fmt.Errorf("service : post : UpdateStatusPostLike : %w", err)
		}
	} else if status == -1 {
		if err := s.repo.UpdateStatusPostLike(0, postLike); err != nil {
			return fmt.Errorf("service : post : UpdateStatusPostLike : %w", err)
		}
	} else {
		if err := s.repo.UpdateStatusPostLike(-1, postLike); err != nil {
			return fmt.Errorf("service : post : UpdateStatusPostLike : %w", err)
		}
	}
	return nil
}

func (s *VotePostService) GetAllLikesByPostId(postId int) (int, error) {
	likes, err := s.repo.GetAllLikesByPostId(postId)
	if err != nil {
		return 0, fmt.Errorf("service : post :  GetAllLikesByPostId : %w", err)
	}
	return likes, nil
}

func (s *VotePostService) GetAllDislikesByPostId(postId int) (int, error) {
	dislikes, err := s.repo.GetAllDislikesByPostId(postId)
	if err != nil {
		return 0, fmt.Errorf("service : post :  GetAllLikesByPostId : %w", err)
	}
	return dislikes, nil
}
