package models

type PostLike struct {
	ID        int
	UserID    int
	CommentID int
	Positive  bool
}
