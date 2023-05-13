package service

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
)

type Notification interface {
	AddNewNotification(notification models.Notification) error
	GetAllNotification(user models.User) ([]models.Notification, error)
	DeleteNotification(notificationId int) error
}

type NotificationService struct {
	repo repository.Notification
}

func NewNotificationService(repo repository.Notification) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) AddNewNotification(notification models.Notification) error {
	err := s.repo.AddNewNotification(notification)
	if err != nil {
		return fmt.Errorf("service AddNewNotification: %w", err)
	}
	return nil
}

func (s *NotificationService) GetAllNotification(user models.User) ([]models.Notification, error) {
	if user == (models.User{}) {
		return nil, nil
	}
	notifies, err := s.repo.GetAllNotification(user.Username)
	if err != nil {
		return nil, fmt.Errorf("service : GetAllNotification %w", err)
	}
	return notifies, nil
}

func (s *NotificationService) DeleteNotification(notificationId int) error {
	err := s.repo.DeleteNotification(notificationId)
	if err != nil {
		return fmt.Errorf("service : delete notification: %w", err)
	}
	return nil
}
