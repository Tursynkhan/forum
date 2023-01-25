package models

type Comment struct {
	ID      int
	Content string
	Author  string
	PostID  int
	UserID  int

	Likes    int
	Dislikes int
}
