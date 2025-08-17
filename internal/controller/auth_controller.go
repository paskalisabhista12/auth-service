package controller

import (
	requestDto "auth-service/internal/model/dto/request"
	responseDto "auth-service/internal/model/dto/response"
	service "auth-service/internal/service"
	"auth-service/pkg/utils"
	exception "auth-service/pkg/utils/exception"
	response "auth-service/pkg/utils/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

type ResponseWrapper struct {
	User any `json:"user"`
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService}
}

func (ac *AuthController) RegisterRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", ac.Register)
		authGroup.POST("/login", ac.Login)
		authGroup.GET("/verify", ac.Verify)
		authGroup.POST("/logout", ac.Logout)
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req requestDto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(exception.ErrBadRequest)
		return
	}

	if err := ac.authService.Register(req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, nil, "User registered successfully")
}

func (ac *AuthController) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(exception.ErrBadRequest)
		return
	}

	token, err := ac.authService.Login(req.Email, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"authToken": token})
}

func (ac *AuthController) Verify(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.Error(exception.ErrBadRequest)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.Error(exception.NewUnauthorizedBusinessException("Invalid authorization header format"))
		return
	}
	token := parts[1]

	data, err := ac.authService.Verify(token)
	if err != nil {
		c.Error(err)
		return
	}

	userResponse, err := utils.UnmarshalDynamic[responseDto.UserResponse]([]byte(data), "user")
	if err != nil {
		panic(err)
	}

	response.Success(c, http.StatusOK, gin.H{"user": userResponse}, "Token is valid")
}

func (ac *AuthController) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.Error(exception.ErrBadRequest)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.Error(exception.NewUnauthorizedBusinessException("Invalid authorization header format"))
		return
	}
	token := parts[1]

	err := ac.authService.Logout(token)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, nil, "Logout success")
}
