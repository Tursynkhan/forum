package repository

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type PostRepository struct {
	db *sql.DB
}
type Post interface {
	CreatePost(post models.Post) error
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post models.Post) error {
	_, err := r.db.Exec("INSERT INTO posts (Title,Content) VALUES (?,?)", post.Title, post.Content)
	if err != nil {
		return fmt.Errorf("repository: create post: %w", err)
	}
	return nil
}
