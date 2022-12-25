package repository

import (
	"database/sql"

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
	usertable = `CREATE TABLE IF NOT EXISTS users(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Username TEXT,
			Email TEXT,
			Password TEXT,
			RetypePassword TEXT,
			Token TEXT
		);`
	postTable = `CREATE TABLE IF NOT EXISTS posts(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT NOT NULL,
			Content TEXT,
			UserId INTEGER,
			FOREIGN KEY (UserId) REFERENCES users (Id) ON DELETE CASCADE
		);`
	commentTable = `CREATE TABLE IF NOT EXISTS comment( 
			Id INTEGER PRIMARY  KEY AUTOINCREMENT,
			Content TEXT,
			UserId INTEGER,
			PostId INTEGER,
			FOREIGN KEY (UserId) REFERENCES users (Id) ON DELETE CASCADE,
			FOREIGN KEY (PostId) REFERENCES posts (Id) ON DELETE CASCADE
		);`
)

func InitDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Username, cfg.DBName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTables(db *sql.DB) error {
	allTables := []string{usertable, postTable, commentTable}
	for _, table := range allTables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}
