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
	CreatePost(post models.Post) (int, error)
	GetAllPosts() ([]models.PostInfo, error)
	GetPost(id int) (models.PostInfo, error)
	CreatePostCategory(postId int, categories []string) error
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post models.Post) (int, error) {
	res, err := r.db.Exec("INSERT INTO posts (Title,Content,UserId) VALUES (?,?,?)", post.Title, post.Content, post.UserID)
	if err != nil {
		return 0, fmt.Errorf("repository : create post : %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *PostRepository) CreatePostCategory(postId int, categories []string) error {
	for _, category := range categories {
		_, err := r.db.Exec("INSERT INTO post_categories (PostId,CategoryId) VALUES (?,?)", postId, category)
		if err != nil {
			return fmt.Errorf("repository : create post : %w", err)
		}
	}
	return nil
}

func (r *PostRepository) GetAllPosts() ([]models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id, users.Username, posts.Title, posts.Content,posts.UserId from posts JOIN users ON users.Id = posts.UserId")
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : get all posts : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		categories_rows, _ := r.db.Query("SELECT categories.Name FROM post_categories JOIN categories ON categories.Id = post_categories.CategoryId WHERE PostId = ?", &p.ID)
		for categories_rows.Next() {
			category := ""
			categories_rows.Scan(&category)
			p.Categories = append(p.Categories, category)
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPost(id int) (models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id, users.Username, posts.Title, posts.Content,posts.UserId from posts JOIN users ON users.Id = posts.UserId WHERE posts.Id=$1", id)
	if err != nil {
		return models.PostInfo{}, fmt.Errorf("repository : get all posts: %w", err)
	}
	var post models.PostInfo
	for rows.Next() {
		err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.UserId)
		if err == sql.ErrNoRows {
			return models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return models.PostInfo{}, err
		}
		categories_rows, _ := r.db.Query("SELECT categories.Name FROM post_categories JOIN categories ON categories.Id = post_categories.CategoryId WHERE PostId = ?", &post.ID)
		for categories_rows.Next() {
			category := ""
			categories_rows.Scan(&category)
			post.Categories = append(post.Categories, category)
		}
	}
	return post, nil
}
