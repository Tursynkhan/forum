package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

type (
	Comment interface {
		CreateComment(comment models.Comment) error
		GetAllComments(postId int) ([]models.Comment, error)
		GetCommentById(id int) (models.Comment, error)
		GetPostById(post models.Comment) (models.Post, error)
	}
)

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(comment models.Comment) error {
	_, err := r.db.Exec("INSERT INTO comments (Content,UserId,PostId) VALUES (?,?,?)", comment.Content, comment.UserID, comment.PostID)
	if err != nil {
		return fmt.Errorf("commentary : create comment : %w", err)
	}

	return nil
}

func (r *CommentRepository) GetAllComments(postId int) ([]models.Comment, error) {
	rows, err := r.db.Query("SELECT comments.Id,comments.Content,comments.UserId,comments.PostId,users.Username FROM comments JOIN users ON users.Id=comments.UserId WHERE comments.PostId=?", postId)
	if err != nil {
		return []models.Comment{}, fmt.Errorf("repository : get all posts : %w", err)
	}
	var comments []models.Comment
	for rows.Next() {
		c := models.Comment{}
		err := rows.Scan(&c.ID, &c.Content, &c.UserID, &c.PostID, &c.Author)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Comment{}, errors.New("No comments")
		} else if err != nil {
			return []models.Comment{}, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *CommentRepository) GetCommentById(id int) (models.Comment, error) {
	row := r.db.QueryRow("SELECT comments.Id,comments.Content,comments.UserId,comments.PostId,users.Username FROM comments JOIN users ON users.Id=comments.UserId  WHERE comments.Id=?", id)
	c := models.Comment{}
	err := row.Scan(&c.ID, &c.Content, &c.UserID, &c.PostID, &c.Author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Comment{}, err
		} else {
			return models.Comment{}, err
		}
	}
	return c, nil
}

func (r *CommentRepository) GetPostById(post models.Comment) (models.Post, error) {
	row := r.db.QueryRow("SELECT Id,Title,Content,UserId,Created FROM posts WHERE Id=?", post.PostID)
	var p models.Post
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &p.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, err
		} else {
			return models.Post{}, fmt.Errorf("repository: GetLenAllPost : %w", err)
		}
	}
	return p, nil
}
