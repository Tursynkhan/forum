package repository

import (
	"database/sql"
	"fmt"
)

type Session interface {
	DeleteSessionById() error
}

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) DeleteSessionById() error {
	rows, err := r.db.Query("SELECT Id FROM session WHERE ExpireTime < datetime('now')")
	if err != nil {
		return fmt.Errorf("repository : DeleteSessionById : %w", err)
	}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("repository : DeleteSessionById %w", err)
		}
		_, err = r.db.Exec("DELETE FROM session WHERE Id =?", id)
		if err != nil {
			return fmt.Errorf("repository : DeleteSessionById %w", err)
		}
	}
	return nil
}
