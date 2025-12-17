package repository

import (
	"auth-service/internal/model"
	"gorm.io/gorm"
)

type EndpointRepository interface {
	FindByServicePathAndHttpMethod(service string, path string, httpMethod string) (model.Endpoint, error)
}

type endpointRepository struct {
	db *gorm.DB
}

func NewEndpointRepository(db *gorm.DB) EndpointRepository {
	return &endpointRepository{db}
}

func (r *endpointRepository) FindByServicePathAndHttpMethod(service string, path string, httpMethod string) (model.Endpoint, error) {
	var endpoint model.Endpoint
	result := r.db.Preload("Permission").Where("service = ? AND path = ? AND http_method = ?", service, path, httpMethod).First(&endpoint)
	return endpoint, result.Error
}
