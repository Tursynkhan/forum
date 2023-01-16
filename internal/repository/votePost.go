package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

type VotePostLikeRepository struct {
	db *sql.DB
}

type VotePost interface {
	CreateLikePost(postLike models.PostLike) error
	GetStatusPostLike(postLike models.PostLike) (int, error)
	UpdateStatusPostLike(status int, postLike models.PostLike) error
}

func NewVotePostRepository(db *sql.DB) *VotePostLikeRepository {
	return &VotePostLikeRepository{db: db}
}

func (r *VotePostLikeRepository) CreateLikePost(postLike models.PostLike) error {
	_, err := r.db.Exec("INSERT INTO posts_like (UserId,PostId,Status) VALUES (?,?,?)", postLike.UserID, postLike.PostID, postLike.Status)
	if err != nil {
		return fmt.Errorf("CreateLikePost : create  likePost : %w", err)
	}
	return nil
}

func (r *VotePostLikeRepository) GetStatusPostLike(postlike models.PostLike) (int, error) {
	row := r.db.QueryRow("SELECT Status FROM posts_like WHERE PostId=? AND UserId=?", postlike.PostID, postlike.UserID)
	status := 0
	err := row.Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, err
		} else {
			return 0, fmt.Errorf("GetStatusPostLike : getStatus : %w", err)
		}
	}
	return status, nil
}

func (r *VotePostLikeRepository) UpdateStatusPostLike(status int, postLike models.PostLike) error {
	query := fmt.Sprintf("UPDATE posts_like SET Status=%d WHERE PostId=? AND UserId=?", status)
	_, err := r.db.Exec(query, postLike.PostID, postLike.UserID)
	if err != nil {
		return fmt.Errorf("repository: updateStatus : %w", err)
	}
	return nil
}
