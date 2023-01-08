package repository

import (
	"database/sql"
)

type (
	Comment interface{}
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
	}
}
