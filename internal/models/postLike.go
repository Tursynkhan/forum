package models

type PostLike struct {
	ID       int
	UserID   int
	PostID   int
	Positive int
}
