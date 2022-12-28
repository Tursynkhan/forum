package repository

import (
	"database/sql"
)

type (
	Post    interface{}
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
	}
}
