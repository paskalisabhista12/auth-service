package model

type Permission struct {
	PermissionID uint   `gorm:"primaryKey;column:permission_id"`
	Name         string `gorm:"column:name"`
	Description  string `gorm:"column:description"`
}
