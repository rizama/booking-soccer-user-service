package repositories

import (
	"context"
	"errors"
	errWrap "user-service/common/error"
	errConst "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Register(context.Context, *dto.RegisterRequest) (*models.User, error)
	Update(context.Context, *dto.UpdateRequest, string) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	FindByUUID(context.Context, string) (*models.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	user := models.User{
		UUID:        uuid.New(),
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		RoleId:      req.RoleId,
	}

	err := ur.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, errWrap.WrapErr(errConst.ErrSQLError)
	}

	return &user, nil
}

func (ur *UserRepository) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*models.User, error) {
	user := models.User{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    *req.Password,
		PhoneNumber: req.PhoneNumber,
	}

	err := ur.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&user).Error
	if err != nil {
		return nil, errWrap.WrapErr(errConst.ErrSQLError)
	}

	return &user, nil
}

func (ur *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User // prepare struct for result
	err := ur.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapErr(errConst.ErrUserNotFound)
		}
		return nil, errWrap.WrapErr(errConst.ErrSQLError)
	}
	return &user, nil
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User // prepare struct for result
	err := ur.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapErr(errConst.ErrUserNotFound)
		}
		return nil, errWrap.WrapErr(errConst.ErrSQLError)
	}
	return &user, nil
}

func (ur *UserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User // prepare struct for result
	err := ur.db.WithContext(ctx).Preload("Role").Where("uuid =?", uuid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapErr(errConst.ErrUserNotFound)
		}
		return nil, errWrap.WrapErr(errConst.ErrSQLError)
	}
	return &user, nil
}
