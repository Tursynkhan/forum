package service

import (
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
	return s.repo.CreateLikePost(postLike)
}
