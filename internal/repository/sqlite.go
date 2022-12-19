package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

const (
	usertable = `CREATE TABLE users(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Username TEXT,
			Password TEXT
		)`
	postTable = `CREATE TABLE posts(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT NOT NULL,
			Content TEXT,
			FOREIGN KEY (UserId) REFERENCES users (Id) ON DELETE CASCADE
		)`
	commentTable = `CREATE TABLE comment(
			Id INTEGER PRIMARY  KEY AUTOINCREMENT,
			Content TEXT,
			FOREIGN KEY (UserId) REFERENCES users (Id) ON DELETE CASCADE
			FOREIGN KEY (PostId) REFERENCES posts (Id) ON DELETE CASCADE
		)`
)

func InitDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
