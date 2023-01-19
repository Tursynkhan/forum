package models

type Info struct {
	User        User
	Post        PostInfo
	Posts       []PostInfo
	Comments    []Comment
	Category    []Category
	PostLike    PostLike
	CommentLike CommentLike
}
