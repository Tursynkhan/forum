package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"net/mail"
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
	userCheck, err := s.repo.GetUser(user.Username)
	if err != nil {
		log.Println(err)
	}
	if userCheck.Username == user.Username {
		return errors.New("User exist")
	}
	if err := checkUser(user); err != nil {
		return err
	}
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, time.Time, error) {
	user, err := s.repo.GetUser(username)
	log.Println(user)
	if err != nil {
		return "", time.Time{}, errors.New("User don't exist")
	}
	if user.Username == "" {
		return "", time.Time{}, errors.New("User don't exist")
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

func checkUser(user models.User) error {
	for _, letter := range user.Username {
		if letter < 32 || letter > 126 {
			return errors.New("Incorrect input")
		}
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.New("Invalid email")
	}
	if len(user.Username) < 2 || len(user.Username) > 36 {
		return errors.New("Invalid username")
	}
	if user.Password != user.RepeatPassword {
		return errors.New("Password doesn't match")
	}
	return nil
}
