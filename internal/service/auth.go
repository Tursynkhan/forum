package service

import (
	"crypto/sha1"
	"fmt"

	"forum/internal/models"
	"forum/internal/repository"
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

// func (s *AuthService) GenerateToken(username, password string) (string, error) {
// 	//get user from DB
// }

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
