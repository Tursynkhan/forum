package repository

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type VotePostLikeRepository struct {
	db *sql.DB
}
type VotePost interface {
	CreateLikePost(postLike models.PostLike) error
}

func NewVotePostRepository(db *sql.DB) *VotePostLikeRepository {
	return &VotePostLikeRepository{db: db}
}

func (r *VotePostLikeRepository) CreateLikePost(postLike models.PostLike) error {
	_, err := r.db.Exec("INSERT INTO posts_like (UserId,PostId,Positive) VALUES (?,?,?)", postLike.UserID, postLike.PostID, postLike.Positive)
	if err != nil {
		return fmt.Errorf("CreateLikePost : create  likePost : %w", err)
	}
	return nil
}
