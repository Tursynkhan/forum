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
	User
	Notification
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Autorization: NewAuthRepository(db),
		Post:         NewPostRepository(db),
		Comment:      NewCommentRepository(db),
		VotePost:     NewVotePostRepository(db),
		VoteComment:  NewVoteCommentRepository(db),
		User:         NewUserRepository(db),
		Notification: NewNotificationRepository(db),
	}
}
