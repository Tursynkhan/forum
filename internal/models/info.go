package models

type Info struct {
	User          User
	ProfileUser   ProfileUser
	Notifications []Notification
	Post          PostInfo
	Posts         []PostInfo
	Comments      []Comment
	Category      []Category
	PostLike      PostLike
	CommentLike   CommentLike
}
