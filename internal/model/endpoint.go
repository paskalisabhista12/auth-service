package model

import ()

type Endpoint struct {
	EndpointID   int    `gorm:"column:endpoint_id;primaryKey;autoIncrement" json:"endpoint_id"`
	Service      string `gorm:"column:service" json:"service"`
	Path         string `gorm:"column:path" json:"path"`
	HTTPMethod   string `gorm:"column:http_method" json:"http_method"`
	PermissionID int    `gorm:"column:permission_id" json:"permission_id"`

	Permission Permission `gorm:"foreignKey:PermissionID;references:PermissionID" json:"permission"`
}
