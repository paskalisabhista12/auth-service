package controller

import (
	requestDto "auth-service/internal/model/dto/request"
	responseDto "auth-service/internal/model/dto/response"
	service "auth-service/internal/service"
	"auth-service/pkg/utils"
	exception "auth-service/pkg/utils/exception"
	response "auth-service/pkg/utils/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
		authGroup.POST("/introspect", ac.Introspect)
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
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"user": userResponse}, "Token is valid")
}

func (ac *AuthController) Introspect(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	var req struct {
		Service  string `json:"service" binding:"required"`
		Endpoint string `json:"endpoint" binding:"required"`
		Method   string `json:"method" binding:"required"`
	}

	if authHeader == "" {
		c.Error(exception.ErrBadRequest)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
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
		c.Error(err)
		return
	}

	/*
		Create user role_permission check
	*/
	err = ac.authService.EnforceAuthorization(userResponse.Email, req.Service, req.Endpoint, req.Method)
	if err != nil {
		c.Error(err)
		return
	}

	userResponseJSON, err := json.Marshal(userResponse)
	if err != nil {
		c.Error(err)
		return
	}

	/*
		Attach user info into X-User headers
	*/
	c.Header("X-User", string(userResponseJSON))

	response.Success(c, http.StatusOK, nil, "Access granted")
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
