package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum/internal/models"
)

type AuthSql struct {
	db *sql.DB
}
type Autorization interface {
	CreateUser(user models.User) error
	GetUser(username string) (models.User, error)
	SaveToken(username, sessionToken string, time time.Time) error
}

func NewAuthRepository(db *sql.DB) *AuthSql {
	return &AuthSql{db: db}
}

func (r *AuthSql) CreateUser(user models.User) error {
	_, err := r.db.Exec("INSERT INTO users (Username,Email,Password) VALUES (?,?,?)", user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("repository: create user: %w", err)
	}
	return nil
}

func (r *AuthSql) GetUser(username string) (models.User, error) {
	rows, err := r.db.Query("SELECT Id,Username,Password from users WHERE username=$1 ", username)
	if err != nil {
		return models.User{}, fmt.Errorf("repository: get user: %w", err)
	}
	var user models.User
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err == sql.ErrNoRows {
			return models.User{}, errors.New("No user with that username")
		} else if err != nil {
			return models.User{}, err
		}
	}
	return user, nil
}

func (r *AuthSql) SaveToken(username, sessionToken string, time time.Time) error {
	_, err := r.db.Exec("UPDATE users SET Token=$1,ExpireTime=$2 WHERE Username=$3", sessionToken, time, username)
	if err != nil {
		return fmt.Errorf("repository: save token: %w", err)
	}
	return nil
}
