package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

type VoteCommentRepository struct {
	db *sql.DB
}

type VoteComment interface {
	CreateLikeComment(comment models.CommentLike) error
	GetStatusCommentLike(comment models.CommentLike) (int, error)
	UpdateStatusCommentLike(status int, comment models.CommentLike) error
	CreateDisLikeComment(comment models.CommentLike) error
	GetCommentLikesByCommentID(id int) (int, error)
	GetCommentDislikesByCommentID(id int) (int, error)
	GetCommentById(comment models.CommentLike) (models.Comment, error)
}

func NewVoteCommentRepository(db *sql.DB) *VoteCommentRepository {
	return &VoteCommentRepository{db: db}
}

func (r *VoteCommentRepository) CreateLikeComment(comment models.CommentLike) error {
	_, err := r.db.Exec("INSERT INTO comments_like (UserId,CommentId,Status) VALUES (?,?,?)", comment.UserID, comment.CommentID, comment.Status)
	if err != nil {
		return fmt.Errorf("CreateLikeComment : %w", err)
	}
	return nil
}

func (r *VoteCommentRepository) GetStatusCommentLike(comment models.CommentLike) (int, error) {
	row := r.db.QueryRow("SELECT Status FROM comments_like WHERE CommentId=? AND UserId=?", comment.CommentID, comment.UserID)
	status := 0
	err := row.Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, err
		} else {
			return 0, fmt.Errorf("GetStatusCommentLike : getStatus : %w", err)
		}
	}
	return status, nil
}

func (r *VoteCommentRepository) UpdateStatusCommentLike(status int, comment models.CommentLike) error {
	query := fmt.Sprintf("UPDATE comments_like SET Status=%d WHERE CommentId=? AND UserId=?", status)
	_, err := r.db.Exec(query, comment.CommentID, comment.UserID)
	if err != nil {
		return fmt.Errorf("repository: updateStatus : %w", err)
	}
	return nil
}

func (r *VoteCommentRepository) CreateDisLikeComment(comment models.CommentLike) error {
	_, err := r.db.Exec("INSERT INTO comments_like (UserId,CommentId,Status) VALUES (?,?,?)", comment.UserID, comment.CommentID, comment.Status)
	if err != nil {
		return fmt.Errorf("CreateLikeComment : %w", err)
	}
	return nil
}

func (r *VoteCommentRepository) GetCommentLikesByCommentID(id int) (int, error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM comments_like WHERE Status=1 AND CommentId=?", id)
	likes := 0
	err := row.Scan(&likes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, err
		} else {
			return 0, fmt.Errorf("GetCommentLikesByCommentID : getlikes : %w", err)
		}
	}
	return likes, nil
}

func (r *VoteCommentRepository) GetCommentDislikesByCommentID(id int) (int, error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM comments_like WHERE Status=-1 AND CommentId=?", id)
	likes := 0
	err := row.Scan(&likes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, err
		} else {
			return 0, fmt.Errorf("GetCommentDislikesByCommentID : getDislikes : %w", err)
		}
	}
	return likes, nil
}

func (r *VoteCommentRepository) GetCommentById(comment models.CommentLike) (models.Comment, error) {
	row := r.db.QueryRow("SELECT comments.Id FROM comments  WHERE comments.Id=?", comment.CommentID)
	var c models.Comment
	err := row.Scan(&c.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Comment{}, err
		} else {
			return models.Comment{}, fmt.Errorf("repository: GetCommentById : %w", err)
		}
	}
	return c, nil
}
