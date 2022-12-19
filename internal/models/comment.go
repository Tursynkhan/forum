package models

type Comment struct {
	ID      int
	PostID  int
	UserID  int
	Content string
}
