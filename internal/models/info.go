package models

type Info struct {
	User          User
	ProfileUser   ProfileUser
	Notifications []Notification
	Post          PostInfo
	Posts         []PostInfo
	Comments      []Comment
	Category      []Category
	Roles         []Role
	PostLike      PostLike
	CommentLike   CommentLike
}
