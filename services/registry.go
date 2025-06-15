package services

import (
	"user-service/repositories"
	services "user-service/services/user"
)

type IServiceRegistry interface {
	GetUserService() services.IUserService
}

type Registry struct {
	repository repositories.IRepositryRegistry
}

func NewServiceRegistry(repository repositories.IRepositryRegistry) IServiceRegistry {
	return &Registry{
		repository: repository,
	}
}

func (r *Registry) GetUserService() services.IUserService {
	return services.NewUserService(r.repository)
}
