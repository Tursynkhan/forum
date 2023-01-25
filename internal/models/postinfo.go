package models

type PostInfo struct {
	ID         int
	Author     string
	Categories []string
	Title      string
	Content    string
	Created    string
	UserId     int
	Likes      int
	Dislikes   int
}
