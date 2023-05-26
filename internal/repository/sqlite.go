package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Host     string
	Port     string
	Name     string
	Password string
	DBName   string
	SSLMode  string
}

const (
	usertable = `CREATE TABLE IF NOT EXISTS users(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Username TEXT UNIQUE,
			Email TEXT UNIQUE,
			Password TEXT,
			RoleId INTEGER,
			FOREIGN KEY (RoleId) REFERENCES roles (Id) ON DELETE CASCADE
			);`
	sessionTable = `CREATE TABLE IF NOT EXISTS session(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Token TEXT,
			ExpireTime DATETIME,
			UserId INTEGER,
			FOREIGN KEY (UserId) REFERENCES users (Id) ON DELETE CASCADE
			);`
	roleTable = `CREATE TABLE IF NOT EXISTS roles(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT UNIQUE
			);`
	insertRoles = `INSERT OR IGNORE INTO roles(Name) VALUES
			('User'),
			('Guest'),
			('Administrator'),
			('Moderator');
			`
	postTable = `CREATE TABLE IF NOT EXISTS posts(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT NOT NULL,
			Content TEXT,
			Created TEXT,
			UserId INTEGER,
			Approved TEXT DEFAULT "NO",
			ReportStatus TEXT,
			FOREIGN KEY (UserId) REFERENCES users (Id) ON DELETE CASCADE
			);`
	categoryTable = `CREATE TABLE IF NOT EXISTS categories(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT UNIQUE
			);`
	postCategoryTable = `CREATE TABLE IF NOT EXISTS post_categories(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			PostId INTEGER,
			CategoryId INTEGER,
			FOREIGN KEY (PostId) REFERENCES posts (Id) ON DELETE CASCADE,
			FOREIGN KEY (CategoryId) REFERENCES categories (Id) ON DELETE CASCADE
			);`
	commentTable = `CREATE TABLE IF NOT EXISTS comments( 
			Id INTEGER PRIMARY  KEY AUTOINCREMENT,
			Content TEXT,
			UserId INTEGER,
			PostId INTEGER,
			FOREIGN KEY (UserId) REFERENCES users (Id) ON DELETE CASCADE,
			FOREIGN KEY (PostId) REFERENCES posts (Id) ON DELETE CASCADE
			);`
	postLikeTable = `CREATE TABLE IF NOT EXISTS posts_like(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			UserId INTEGER,
			PostId INTEGER,
			Status INTEGER DEFAULT 0,
			FOREIGN KEY (UserId) REFERENCES users (Id) ON DELETE CASCADE,
			FOREIGN KEY (PostId) REFERENCES posts (Id) ON DELETE CASCADE
			);`
	commentLikeTable = `CREATE TABLE IF NOT EXISTS comments_like(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			UserId INTEGER,
			CommentId INTEGER,
			Status INTEGER DEFAULT 0,
			FOREIGN KEY (UserId) REFERENCES users (Id) ON DELETE CASCADE,
			FOREIGN KEY (CommentId) REFERENCES comments (Id) ON DELETE CASCADE
			);`
	imageTable = `CREATE TABLE IF NOT EXISTS images(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			ImageName TEXT,
			PostId INTEGER,
			FOREIGN KEY (PostId) REFERENCES posts (Id) ON DELETE CASCADE
			);`
	notificationTable = `CREATE TABLE IF NOT EXISTS notifications(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			FromUser TEXT,
			ToUser TEXT,
			Content TEXT,
			PostId INTEGER,
			TimeStamp TEXT,
			IsRead INTEGER,
			FOREIGN KEY (PostId) REFERENCES posts (Id) ON DELETE CASCADE
			);`
	insertCategories = `INSERT OR IGNORE INTO categories(Name) VALUES
			('Getting Help'),
			('Releases'),
			('Technical Discussion'),
			('Community'),
			('Jobs'),
			('Site Feedback');
			`
)

func InitDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Name, cfg.DBName+"?_foreign_keys=on")
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
	allTables := []string{usertable, postTable, commentTable, categoryTable, postCategoryTable, postLikeTable, commentLikeTable, insertCategories, sessionTable, imageTable, notificationTable, roleTable, insertRoles}
	for _, table := range allTables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}
