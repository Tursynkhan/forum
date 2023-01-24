package service

import (
	"forum/internal/repository"
)

type Session interface {
	CompareExpirationTime() error
}

type SessionService struct {
	repo repository.Session
}

func NewSessionService(repo repository.Session) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) CompareExpirationTime() error {
	return s.repo.DeleteSessionById()
}
