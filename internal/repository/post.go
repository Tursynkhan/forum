package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/models"
)

type PostRepository struct {
	db *sql.DB
}
type Post interface {
	CreatePost(post models.Post) error
	GetAllPosts() ([]models.Post, error)
	GetPost(id int) (models.Post, error)
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

func (r *PostRepository) GetAllPosts() ([]models.Post, error) {
	rows, err := r.db.Query("SELECT Id,Title,Content from posts")
	if err != nil {
		return []models.Post{}, fmt.Errorf("repository: get all posts: %w", err)
	}
	var posts []models.Post
	for rows.Next() {
		p := models.Post{}
		err := rows.Scan(&p.ID, &p.Title, &p.Content)
		if err == sql.ErrNoRows {
			return []models.Post{}, errors.New("No posts")
		} else if err != nil {
			return []models.Post{}, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPost(id int) (models.Post, error) {
	rows, err := r.db.Query("SELECT Id,Title,Content from posts WHERE Id=$1", id)
	if err != nil {
		return models.Post{}, fmt.Errorf("repository: get all posts: %w", err)
	}
	var post models.Post
	for rows.Next() {
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err == sql.ErrNoRows {
			return models.Post{}, errors.New("No posts")
		} else if err != nil {
			return models.Post{}, err
		}
	}
	return post, nil
}
