package models

type User struct {
	ID             int
	Username       string
	Email          string
	Password       string
	RepeatPassword string

	Token string
}
