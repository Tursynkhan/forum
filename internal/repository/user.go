package repository

import (
	"database/sql"
	"forum/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

type User interface {
	GetPostByUsername(username string) ([]models.PostInfo, error)
	GetLikedPostByUsername(usernaem string) ([]models.PostInfo, error)
	GetCommentPostByUsername(username string) ([]models.PostInfo, error)
	GetUserByUsername(username string) (models.User, error)
}

func (r *UserRepository) GetPostByUsername(username string) ([]models.PostInfo, error) {
}

func (r *UserRepository) GetLikedPostByUsername(usernaem string) ([]models.PostInfo, error) {
}

func (r *UserRepository) GetCommentPostByUsername(username string) ([]models.PostInfo, error) {
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
}
