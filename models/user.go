package models

import (
	"time"
)

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"type:varchar(100);not null" json:"name"`
    Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
    Password  string    `gorm:"type:varchar(255);not null" json:"-"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}