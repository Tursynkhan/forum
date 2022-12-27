package repository

import (
	"database/sql"
	"log"
	"time"

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

func (r *AuthSql) GetUser(username string) (models.User, error) {
	rows, err := r.db.Query("SELECT Id,Username,Password from users WHERE username=? ", username)
	if err != nil {
		return models.User{}, err
	}
	var user models.User
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err == sql.ErrNoRows {
			return models.User{}, nil
		} else {
			log.Println(err)
		}
	}
	if err = rows.Err(); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *AuthSql) SaveToken(username, sessionToken string, time time.Time) error {
	_, err := r.db.Exec("UPDATE users SET Token=$1,ExpireTime=$2 WHERE Username=$3", sessionToken, time, username)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
