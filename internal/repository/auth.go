package repository

import (
	"database/sql"
	"forum/internal/models"
	"log"
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

func (r *AuthSql) GetUserByName(username, password string) (models.User, error) {
	// var user models.User
	// query:=fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password=$2",usertable)
	// err:=r.db.Get(&user,query,username,password)
	// return user,err
}
