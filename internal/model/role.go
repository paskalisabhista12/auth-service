package model

type Role struct {
	RoleID      uint   `gorm:"primaryKey;column:role_id"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`

	Permissions []Permission `gorm:"many2many:role_permissions;joinForeignKey:RoleID;joinReferences:PermissionID"`
}
