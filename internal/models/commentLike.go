package models

type CommentLike struct {
	ID        int
	UserID    int
	CommentID int
	Status    int
	Likes     int
	Dislikes  int
}
