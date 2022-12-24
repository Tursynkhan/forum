package repository

import (
	"database/sql"

	"forum/internal/models"
)

type (
	Autorization interface {
		CreateUser(user models.User) error
	}
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
