package repository

import (
	"database/sql"

	"forum/internal/models"
)

type AuthSql struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthSql {
	return &AuthSql{db: db}
}

func (r *AuthSql) CreateUser(user models.User) (int, error) {
	return 0, nil
}
