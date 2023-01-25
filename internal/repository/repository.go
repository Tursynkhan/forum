package repository

import (
	"database/sql"
)

type Repository struct {
	Autorization
	Post
	Comment
	VotePost
	VoteComment
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Autorization: NewAuthRepository(db),
		Post:         NewPostRepository(db),
		Comment:      NewCommentRepository(db),
		VotePost:     NewVotePostRepository(db),
		VoteComment:  NewVoteCommentRepository(db),
	}
}
