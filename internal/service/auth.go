package service

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("Invalid email")
	ErrInvalidUsername = errors.New("Invalid username")
	ErrInvalidPassword = errors.New("Invalid password")
	ErrUserNotFound    = errors.New("User not found")
	ErrUserExist       = errors.New("User exist")
)

type Autorization interface {
	CreateUser(user models.User) error
	GenerateToken(username, password string) (string, time.Time, error)
	ParseToken(token string) (models.User, error)
	DeleteToken(token string) error
}

type AuthService struct {
	repo repository.Autorization
}

func NewAuthService(repo repository.Autorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) error {
	userCheck, err := s.repo.GetUser(user.Username)
	if err != nil {
		return err
	}
	if userCheck.Username == user.Username {
		return ErrUserExist
	}
	if err := checkUser(user); err != nil {
		return err
	}
	if user.Password, err = generatePasswordHash(user.Password); err != nil {
		return err
	}
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, time.Time, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", time.Time{}, err
	}
	if user.Username == "" {
		return "", time.Time{}, ErrUserNotFound
	}
	if err := checkHash(user.Password, password); err != nil {
		return "", time.Time{}, fmt.Errorf("service : compare hash and password : %v: %w", err, ErrUserNotFound)
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	if err := s.repo.SaveToken(user.Username, sessionToken, expiresAt); err != nil {
		return "", time.Time{}, err
	}
	return sessionToken, expiresAt, nil
}

func (s *AuthService) ParseToken(token string) (models.User, error) {
	user, err := s.repo.GetUserByToken(token)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *AuthService) DeleteToken(token string) error {
	return s.repo.DeleteToken(token)
}

func generatePasswordHash(password string) (string, error) {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		return "", fmt.Errorf("service: generatePassword: %v", err)
	}
	return string(hash), nil
}

func checkHash(hpass, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hpass), []byte(password))
}

func checkUser(user models.User) error {
	for _, letter := range user.Username {
		if letter < 32 || letter > 126 {
			return fmt.Errorf("service: CreateUser: checkUser: %v", ErrInvalidUsername)
		}
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return fmt.Errorf("service: CreateUser: checkUser: %v", ErrInvalidEmail)
	}
	if len(user.Username) < 2 || len(user.Username) > 36 {
		return fmt.Errorf("service: CreateUser: checkUser: %v", ErrInvalidUsername)
	}
	if user.Password != user.RepeatPassword {
		return fmt.Errorf("service: CreateUser: checUser: %v", errors.New("Password doesn't match"))
	}
	return nil
}
