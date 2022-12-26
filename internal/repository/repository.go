package repository

import (
	"database/sql"
	"time"

	"forum/internal/models"
)

type (
	Autorization interface {
		CreateUser(user models.User) error
		GetUser(username string) (models.User, error)
		SaveToken(username, sessionToken string, time time.Time) error
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
