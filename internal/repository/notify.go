package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

type Notification interface {
	AddNewNotification(notification models.Notification) error
	GetAllNotification(username string) ([]models.Notification, error)
	DeleteNotification(notificationId int) error
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) AddNewNotification(notification models.Notification) error {
	_, err := r.db.Exec("INSERT INTO notifications (FromUser,ToUser,Content,PostId,TimeStamp) VALUES (?,?,?,?)", notification.From, notification.To, notification.Content, notification.PostId, notification.TimeStamp)
	if err != nil {
		return fmt.Errorf("repository: AddNewNotification : %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetAllNotification(username string) ([]models.Notification, error) {
}

func (r *NotificationRepository) DeleteNotification(notificationId int) error {
}
