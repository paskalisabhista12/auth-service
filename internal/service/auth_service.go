package service

import (
	"auth-service/internal/config"
	"auth-service/internal/infra/redis"
	"auth-service/internal/model"
	"auth-service/internal/model/dto/request"
	"auth-service/internal/repository"
	"auth-service/pkg/utils/exception"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(c *gin.Context,req requestDTO.RegisterRequest) error
	Login(c *gin.Context, email, password string) (string, error)
	Verify(c *gin.Context, authToken string) (string, error)
	Logout(c *gin.Context, authToken string) error
	EnforceAuthorization(c *gin.Context, userEmail string, service string, endpoint string, httpMethod string) error
}

type authService struct {
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	endpointRepo repository.EndpointRepository
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, endpointRepo repository.EndpointRepository) AuthService {
	return &authService{userRepo, roleRepo, endpointRepo}
}

func (s *authService) Register(c *gin.Context, req requestDTO.RegisterRequest) error {
	user := model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

	// Check if user already exists
	_, err := s.userRepo.FindByEmail(strings.TrimSpace(user.Email))
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
	if _, err := s.userRepo.Create(user); err != nil {
		return exception.NewInternal("Failed to save user")
	}

	return err

}

func (s *authService) Login(c *gin.Context, email, password string) (string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(strings.TrimSpace(email))
	if err != nil || user.ID == 0 {
		return "", exception.NewUnauthorizedBusinessException("Invalid email or password")
	}

	var roleNames []string
	for _, r := range user.Roles {
		roleNames = append(roleNames, r.Name)
	}

	roleNamesString := strings.Join(roleNames, "|") // e.g. "SUPERADMIN|ADMIN|etc"

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
		"roles":      roleNamesString,
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
			"roles":      roleNamesString,
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

func (s *authService) Verify(c *gin.Context, authToken string) (string, error) {
	authToken = strings.TrimSpace(authToken)
	if authToken == "" {
		return "", exception.NewUnauthorizedBusinessException("Authorization token is required")
	}

	secret := config.LoadConfig().JwtSecret
	_, err := verifyToken(authToken, secret)
	if err != nil {
		return "", err
	}

	data, err := redis.Rdb.Get(redis.Ctx, authToken).Result()
	if err != nil {
		return "", exception.NewUnauthorizedBusinessException("Token not valid or expired")
	}

	return data, nil
}

func (s *authService) Logout(c *gin.Context, authToken string) error {
	authToken = strings.TrimSpace(authToken)
	if authToken == "" {
		return exception.NewUnauthorizedBusinessException("Authorization token is required")
	}

	secret := config.LoadConfig().JwtSecret
	_, err := verifyToken(authToken, secret)
	if err != nil {
		return err
	}

	deleted, err := redis.Rdb.Del(redis.Ctx, authToken).Result()
	if err != nil {
		return exception.NewInternal("Failed to delete token")
	}

	if deleted == 0 {
		return exception.NewNotFound("Token not found")
	}

	return nil
}

func (s *authService) EnforceAuthorization(c *gin.Context, userEmail string, service string, path string, httpMethod string) error {
	/*
		Get user roles
	*/
	user, err := s.userRepo.FindByEmail(userEmail)
	if err != nil {
		return exception.NewUnauthorizedBusinessException("Invalid email or password")
	}
	roleIds := extractRoleIDs(user.Roles)

	/*
		Check endpoint in DB and extract the needed permission to access the endpoint
	*/
	endpoint, err := s.endpointRepo.FindByServicePathAndHttpMethod(service, path, httpMethod)
	if err != nil {
		return exception.NewNotFound("Endpoint not found")
	}
	requiredPermission := endpoint.Permission

	/*
		Extract permissions from roles
	*/
	permissions, _ := s.roleRepo.GetPermissionsByRoleIds(roleIds)

	/*
		Verify endpoint permission to the user privilege
	*/
	isAllowed := containsPermission(permissions, requiredPermission)
	if !isAllowed {
		return exception.NewUnauthorizedBusinessException("User has no permission to access this endpoint")
	}
	return nil
}

func verifyToken(tokenString string, secret string) (jwt.MapClaims, error) {
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

func extractRoleIDs(roles []model.Role) []int {
	ids := make([]int, len(roles))
	for i, r := range roles {
		ids[i] = int(r.RoleID)
	}
	return ids
}

func containsPermission(permissions []model.Permission, permission model.Permission) bool {
	for _, item := range permissions {
		if item.Name == "ALL" { // By default, every endpoint is bypassed if the user has the 'ALL' permission (which should come from the SUPERADMIN role).
			return true
		}
		if item.PermissionID == permission.PermissionID {
			return true
		}
	}
	return false
}
