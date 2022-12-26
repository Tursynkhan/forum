package models

import "time"

type User struct {
	ID             int
	Username       string
	Email          string
	Password       string
	RepeatPassword string

	Token      string
	Expiretime time.Time
}
