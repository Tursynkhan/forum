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
	rows, err := r.db.Query("SELECT Id, Content, UserId,PostId from comments WHERE PostId=?", postId)
	if err != nil {
		return []models.Comment{}, fmt.Errorf("repository : get all posts : %w", err)
	}
	var comments []models.Comment
	for rows.Next() {
		c := models.Comment{}
		err := rows.Scan(&c.ID, &c.Content, &c.UserID, &c.PostID)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Comment{}, errors.New("No comments")
		} else if err != nil {
			return []models.Comment{}, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
