package repository

import (
	"database/sql"
)

type Repository struct {
	Autorization
	Post
	Comment
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Autorization: NewAuthRepository(db),
		Post:         NewPostRepository(db),
		Comment:      NewCommentRepository(db),
	}
}
