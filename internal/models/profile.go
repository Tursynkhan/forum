package models

type ProfileUser struct {
	Username string
	Email    string

	CountOfPosts    int
	CountOfLikes    int
	CountOfComments int
}
