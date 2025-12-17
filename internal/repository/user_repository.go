package repository

import (
	"auth-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]model.User, error)
	Create(user model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *userRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	result := r.db.Preload("Roles").Where("email = ?", email).First(&user)
	return user, result.Error
}

func (r *userRepository) Create(user model.User) (model.User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}
