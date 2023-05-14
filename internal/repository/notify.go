package repository

import (
	"database/sql"
	"errors"
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
	_, err := r.db.Exec("INSERT INTO notifications (FromUser,ToUser,Content,PostId,TimeStamp,IsRead) VALUES (?,?,?,?,?,?)", notification.From, notification.To, notification.Content, notification.PostId, notification.TimeStamp, notification.IsRead)
	if err != nil {
		return fmt.Errorf("repository: AddNewNotification : %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetAllNotification(username string) ([]models.Notification, error) {
	rows, err := r.db.Query("SELECT Id,FromUser,ToUser,Content,PostId,TimeStamp,IsRead FROM notifications WHERE IsRead=0 AND ToUSer=?", username)
	if err != nil {
		return nil, fmt.Errorf("repository : GetAllNotification : %w", err)
	}
	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(&notification.Id, &notification.From, &notification.To, &notification.Content, &notification.PostId, &notification.TimeStamp, &notification.IsRead)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else if err != nil {
			return nil, fmt.Errorf("repository : GetAllNotifications %w", err)
		}

		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (r *NotificationRepository) DeleteNotification(notificationId int) error {
	_, err := r.db.Exec("DELETE FROM notifications WHERE Id=?", notificationId)
	if err != nil {
		return fmt.Errorf("repository : delete notification: %w", err)
	}
	return nil
}
