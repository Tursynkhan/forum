package models

type Info struct {
	User     User
	Post     PostInfo
	Posts    []PostInfo
	Comments []Comment
}
