package models

type CommentLike struct {
	UserID    int
	CommentID int
	Positive  bool
}
