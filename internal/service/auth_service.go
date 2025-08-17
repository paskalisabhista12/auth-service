package service

import (
	config "auth-service/internal/config"
	"auth-service/internal/infra/redis"
	model "auth-service/internal/model"
	requestDTO "auth-service/internal/model/dto/request"
	repository "auth-service/internal/repository"
	exception "auth-service/pkg/utils/exception"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req requestDTO.RegisterRequest) error
	Login(email, password string) (string, error)
	Verify(authToken string) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(req requestDTO.RegisterRequest) error {
	user := model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

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
		return exception.NewInternal("Failed to save user")
	}

	return err

}

func (s *authService) Login(email, password string) (string, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(strings.TrimSpace(email))
	if err != nil || user.ID == 0 {
		return "", exception.NewUnauthorizedBusinessException("Invalid email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", exception.NewUnauthorizedBusinessException("Invalid email or password")
	}

	// Load secret
	secret := config.LoadConfig().JwtSecret
	if secret == "" {
		return "", exception.NewInternal("JWT secret is not set")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"exp":        time.Now().Add(12 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", exception.NewInternal("Failed to sign token")
	}

	value := gin.H{
		"user": gin.H{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		}}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", exception.ErrInternal
	}

	if err := redis.Set(signed, string(jsonValue), 12*3600*time.Second); err != nil {
		return "", exception.ErrInternal
	}

	return signed, nil
}

func (s *authService) Verify(authToken string) (string, error) {
	authToken = strings.TrimSpace(authToken)
	if authToken == "" {
		return "", exception.NewUnauthorizedBusinessException("Authorization token is required")
	}

	secret := config.LoadConfig().JwtSecret
	_, err := VerifyToken(authToken, secret)
	if err != nil {
		return "", err
	}

	data, err := redis.Rdb.Get(redis.Ctx, authToken).Result()
	if err != nil {
		return "", exception.NewUnauthorizedBusinessException("Token not valid or expired")
	}

	return data, nil
}

func VerifyToken(tokenString string, secret string) (jwt.MapClaims, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, exception.NewUnauthorizedBusinessException(
				fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]),
			)
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return nil, exception.NewUnauthorizedBusinessException("Token expired")
			}
		}
		return claims, nil
	}

	return nil, exception.NewUnauthorizedBusinessException("Token invalid")
}
