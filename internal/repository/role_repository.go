package repository

import (
	model "auth-service/internal/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	GetPermissionsByRoleIds(ids []int) ([]model.Permission, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) GetPermissionsByRoleIds(ids []int) ([]model.Permission, error) {
    var roles []model.Role
    result := r.db.Preload("Permissions").Where("role_id IN (?)", ids).Find(&roles)
    if result.Error != nil {
        return nil, result.Error
    }

    // Flatten permissions
    var permissions []model.Permission
    for _, role := range roles {
        permissions = append(permissions, role.Permissions...)
    }

    return permissions, nil
}

