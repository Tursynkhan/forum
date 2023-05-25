package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"time"
)

type AuthSql struct {
	db *sql.DB
}

type Autorization interface {
	CreateUser(user models.User) error
	GetUser(username string) (models.User, error)
	GetEmail(email string) (models.User, error)
	SaveToken(user models.User, sessionToken string, time time.Time) error
	GetUserByToken(token string) (models.User, error)
	GetToken(username string) (string, error)
	UpdateToken(tokenName, newToken string, expireTime time.Time) error
	DeleteToken(token string) error
	DeleteTokenWhenExpireTime() error
}

func NewAuthRepository(db *sql.DB) *AuthSql {
	return &AuthSql{db: db}
}

func (r *AuthSql) CreateUser(user models.User) error {
	_, err := r.db.Exec("INSERT INTO users (Username,Email,Password,RoleId) VALUES (?,?,?,?)", user.Username, user.Email, user.Password, user.RoleID)
	if err != nil {
		return fmt.Errorf("repository : create user : %w", err)
	}
	return nil
}

func (r *AuthSql) GetUser(username string) (models.User, error) {
	rows, err := r.db.Query("SELECT Id,Username,Password from users WHERE username=$1 ", username)
	if err != nil {
		return models.User{}, fmt.Errorf("repository : get user : %w", err)
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

func (r *AuthSql) GetEmail(email string) (models.User, error) {
	row := r.db.QueryRow("SELECT Id,Username,Password,Email FROM users WHERE email=?", email)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, err
		} else {
			return models.User{}, err
		}
	}
	return user, nil
}

func (r *AuthSql) SaveToken(user models.User, sessionToken string, time time.Time) error {
	_, err := r.db.Exec("INSERT INTO session (Token,ExpireTime,UserId) VALUES (?,?,?)", sessionToken, time, user.ID)
	if err != nil {
		return fmt.Errorf("repository : save token : %w", err)
	}
	return nil
}

func (r *AuthSql) GetUserByToken(token string) (models.User, error) {
	row := r.db.QueryRow("SELECT users.Id,users.Username,users.Email,session.ExpireTime,users.RoleId FROM users JOIN session ON users.Id=session.UserId WHERE session.Token=?", token)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Expiretime, &user.RoleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("repository : GetUserByToken : %w", err)
		} else {
			return models.User{}, fmt.Errorf("repository : GetUserByToken : %w", err)
		}
	}
	return user, nil
}

func (r *AuthSql) DeleteToken(token string) error {
	_, err := r.db.Exec("DELETE FROM session WHERE Token=?", token)
	if err != nil {
		return fmt.Errorf("repository : delete token : %w", err)
	}
	return nil
}

func (r *AuthSql) DeleteTokenWhenExpireTime() error {
	_, err := r.db.Exec("DELETE FROM session WHERE ExpireTime <= ?", time.Now())
	if err != nil {
		return fmt.Errorf("repository : DeleteTokenWhenExpireTime : %w", err)
	}
	return nil
}

func (r *AuthSql) UpdateToken(tokenName, newToken string, expireTime time.Time) error {
	_, err := r.db.Exec("UPDATE session SET Token=?, ExpireTime=? WHERE Token=?", newToken, expireTime, tokenName)
	if err != nil {
		return fmt.Errorf("repository : UpdateToken : %w", err)
	}
	return nil
}

func (r *AuthSql) GetToken(username string) (string, error) {
	row := r.db.QueryRow("SELECT session.Token FROM session JOIN users ON users.Id=session.UserId WHERE users.Username=?", username)
	token := ""
	err := row.Scan(&token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", err
		} else {
			return "", fmt.Errorf("repository : GetToken : %w", err)
		}
	}
	return token, nil
}
