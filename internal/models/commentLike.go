package models

type CommentLike struct {
	ID        int
	UserID    int
	CommentID int
	Positive  int
}
