package service

import (
	"forum/internal/repository"
)

type Service struct {
	Autorization
	Post
	Comment
	VotePost
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Autorization: NewAuthService(repos.Autorization),
		Post:         NewPostService(repos.Post),
		Comment:      NewCommentService(repos.Comment),
		VotePost:     NewVotePostService(repos.VotePost),
	}
}
