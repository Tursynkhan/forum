package repository

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type AuthSql struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthSql {
	return &AuthSql{db: db}
}

func (r *AuthSql) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username,email,password) VALUES ($1,$2,$3)")
	row := r.db.QueryRow(query, user.Username, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, nil
	}
	return id, nil
}
