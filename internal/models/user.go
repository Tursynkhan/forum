package models

import "time"

type User struct {
	ID             int
	Username       string
	Email          string
	Password       string
	RepeatPassword string
	RoleID         int
	Role           string
	Token          string
	Expiretime     time.Time
}

type Role struct {
	ID   int
	Name string
}
