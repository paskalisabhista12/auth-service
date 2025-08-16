package service

import (
	config "auth-service/internal/config"
	model "auth-service/internal/model"
	repository "auth-service/internal/repository"
	exception "auth-service/pkg/utils/exception"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user model.User) error
	Login(email, password string) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(user model.User) error {
	// Check if user already exists
	_, err := s.repo.FindByEmail(strings.TrimSpace(user.Email))
	if err == nil {
		return exception.NewConflictBusinessException("User already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return exception.ErrInternal
	}

	user.Password = string(hashedPassword)

	// Save user
	if _, err := s.repo.Create(user); err != nil {
		return exception.NewInternal("failed to save user")
	}

	return err

}

func (s *authService) Login(email, password string) (string, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(strings.TrimSpace(email))
	if err != nil || user.ID == 0 {
		// Use Unauthorized instead of Internal
		return "", exception.NewUnauthorizedBusinessException("invalid email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", exception.NewUnauthorizedBusinessException("invalid email or password")
	}

	// Load secret
	secret := config.LoadConfig().JwtSecret
	if secret == "" {
		return "", exception.NewInternal("JWT secret is not set")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", exception.NewInternal("failed to sign token")
	}

	return signed, nil
}
