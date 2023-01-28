package service

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"net/mail"
	"regexp"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("Invalid email")
	ErrInvalidUsername = errors.New("Invalid username")
	ErrInvalidPassword = errors.New("Invalid password")
	ErrUserNotFound    = errors.New("User not found")
	ErrUserExist       = errors.New("User exist")
	ErrEmailExist      = errors.New("Email exist")
	ErrPasswdNotMatch  = errors.New("Password doesn't match")
)

type Autorization interface {
	CreateUser(user models.User) error
	GenerateToken(username, password string) (string, time.Time, error)
	ParseToken(token string) (models.User, error)
	DeleteToken(token string) error
	DeleteTokenWhenExpireTime() error
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
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("service: create user: %w", err)
		}
	}
	emailCheck, err := s.repo.GetEmail(user.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("service: create user: %w", err)
		}
	}
	if userCheck.Username == user.Username {
		return ErrUserExist
	}
	if emailCheck.Email == user.Email {
		return ErrEmailExist
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
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, ErrUserNotFound
		}
		return "", time.Time{}, err
	}
	if user.Username == "" {
		return "", time.Time{}, ErrUserNotFound
	}
	if err := checkHash(user.Password, password); err != nil {
		return "", time.Time{}, fmt.Errorf("service : compare hash and password : %v: %w", err, ErrUserNotFound)
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(15 * time.Minute)

	if err := s.repo.SaveToken(user, sessionToken, expiresAt); err != nil {
		return "", time.Time{}, err
	}
	return sessionToken, expiresAt, nil
}

func (s *AuthService) ParseToken(token string) (models.User, error) {
	user, err := s.repo.GetUserByToken(token)
	if err != nil {
		return models.User{}, fmt.Errorf("service : ParseToken %w", err)
	}
	return user, nil
}

func (s *AuthService) DeleteToken(token string) error {
	return s.repo.DeleteToken(token)
}

func (s *AuthService) DeleteTokenWhenExpireTime() error {
	return s.repo.DeleteTokenWhenExpireTime()
}

func generatePasswordHash(password string) (string, error) {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		return "", fmt.Errorf("service : generatePassword : %v", err)
	}
	return string(hash), nil
}

func checkHash(hpass, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hpass), []byte(password))
}

func checkUser(user models.User) error {
	for _, letter := range user.Username {
		if letter < 32 || letter > 126 {
			return fmt.Errorf("service: CreateUser: checkUser : %w", ErrInvalidUsername)
		}
	}
	validUsername, err := regexp.MatchString(`[a-zA-Z0-9]{3,12}$`, user.Username)
	if err != nil {
		return err
	}
	if !validUsername {
		return ErrInvalidUsername
	}
	validEmail, err := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, user.Email)
	if err != nil {
		return err
	}
	if !validEmail {
		return ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return fmt.Errorf("service: CreateUser: checkUser : %w", ErrInvalidEmail)
	}
	if len(user.Username) < 2 || len(user.Username) > 36 {
		return fmt.Errorf("service: CreateUser: checkUser : %w", ErrInvalidUsername)
	}

	if !isValidPassword(user.Password) {
		return fmt.Errorf("service : CreateUser : checUser : %w", ErrInvalidPassword)
	}
	if user.Password != user.RepeatPassword {
		return fmt.Errorf("service : CreateUser : checUser : %w", ErrPasswdNotMatch)
	}
	return nil
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	var hasUpper bool
	var hasLower bool
	var hasDigit bool
	var hasSpecial bool
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}
