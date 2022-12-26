package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"forum/internal/models"
	"forum/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const salt = "hkhasdfa2454654asdf1asdf4a5sdf"

type AuthService struct {
	repo repository.Autorization
}

func NewAuthService(repo repository.Autorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) error {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, time.Time, error) {
	// get user from DB
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", time.Time{}, errors.New("User dont exist")
	}
	// if err := checkHash(user.Password, password); err != nil {
	// 	return "", time.Time{}, errors.New("Password doesn't match")
	// }

	sessioToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	if err := s.repo.SaveToken(user.Username, sessioToken, expiresAt); err != nil {
		return "", time.Time{}, err
	}
	return sessioToken, expiresAt, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func checkHash(hpass, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hpass), []byte(password))
}
