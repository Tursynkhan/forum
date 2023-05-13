package models

type Notification struct {
	Id          int
	From        string
	To          string
	Content string
	PostId      int
	CommentId   int
	TimeStamp   string
	IsRead      int
}
