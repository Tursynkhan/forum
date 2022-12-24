package repository

import (
	"database/sql"
	"log"

	"forum/internal/models"
)

type AuthSql struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthSql {
	return &AuthSql{db: db}
}

func (r *AuthSql) CreateUser(user models.User) error {
	_, err := r.db.Exec("INSERT INTO users (Username,Email,Password) VALUES (?,?,?)", user.Username, user.Email, user.Password)
	if err != nil {
		log.Println(err)
	}
	return nil
}
