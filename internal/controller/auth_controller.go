package controller

import (
	model "auth-service/internal/model"
	service "auth-service/internal/service"
	exception "auth-service/pkg/utils/exception"
	response "auth-service/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService}
}

func (ac *AuthController) RegisterRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", ac.Register)
		authGroup.POST("/login", ac.Login)
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req model.User

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

	response.Success(c, http.StatusOK, gin.H{"token": token})
}
