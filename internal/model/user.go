package model

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Email     string    `gorm:"column:email;unique"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	Roles []Role `gorm:"many2many:user_roles;joinForeignKey:UserId;joinReferences:RoleID"`
}
