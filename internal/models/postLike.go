package models

type PostLike struct {
	ID       int
	UserID   int
	PostID   int
	Status   int
	Likes    int
	Dislikes int
}
