package repositories

import (
	repositories "user-service/repositories/user"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositryRegistry interface {
	UserRepo() repositories.IUserRepository
}

// Provider
func NewRepositoryRegistry(db *gorm.DB) IRepositryRegistry {
	return &Registry{db: db}
}

func (r *Registry) UserRepo() repositories.IUserRepository {
	return repositories.NewUserRepository(r.db)
}
