package services

import (
	"context"
	"strings"
	"time"
	"user-service/config"
	"user-service/constants"
	errConst "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// create struct
type UserService struct {
	repository repositories.IRepositryRegistry
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

// create interface
type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(context.Context, *dto.UpdateRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

// create provider function
func NewUserService(repository repositories.IRepositryRegistry) IUserService {
	return &UserService{repository: repository}
}

// create service method of stuct
func (us *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// get existirng user
	user, err := us.repository.UserRepo().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	// generate token
	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute)
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        strings.ToLower(user.Role.Code),
	}

	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime.Unix(), 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	// return token
	return response, nil
}

func (us *UserService) isUsernameExist(ctx context.Context, username string) bool {
	user, err := us.repository.UserRepo().FindByUsername(ctx, username)
	if err != nil {
		return false
	}
	return user != nil
}

func (us *UserService) isEmailExist(ctx context.Context, username string) bool {
	user, err := us.repository.UserRepo().FindByEmail(ctx, username)
	if err != nil {
		return false
	}
	return user != nil
}

func (us *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if us.isUsernameExist(ctx, req.Username) {
		return nil, errConst.ErrUsernameExist
	}

	if us.isEmailExist(ctx, req.Email) {
		return nil, errConst.ErrEmailExist
	}

	if req.Password != req.ConfirmPassword {
		return nil, errConst.ErrPasswordDoesNotMatch
	}

	user, err := us.repository.UserRepo().Register(ctx, &dto.RegisterRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    string(hashPassword),
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		RoleId:      constants.Customer,
	})
	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
		},
	}

	return response, nil
}

func (us *UserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	var (
		password                  string
		checkUsername, checkEmail *models.User
		hashedPassword            []byte
		user, userResult          *models.User
		err                       error
		data                      dto.UserResponse
	)

	user, err = us.repository.UserRepo().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	isUsernameExist := us.isUsernameExist(ctx, req.Username)
	if isUsernameExist && user.Username != req.Username {
		checkUsername, err = us.repository.UserRepo().FindByUsername(ctx, req.Username)
		if err != nil {
			return nil, err
		}

		if checkUsername != nil {
			return nil, errConst.ErrUsernameExist
		}
	}

	isEmailExist := us.isEmailExist(ctx, req.Email)
	if isEmailExist && user.Email != req.Email {
		checkEmail, err = us.repository.UserRepo().FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}

		if checkEmail != nil {
			return nil, errConst.ErrEmailExist
		}
	}

	if req.Password != nil {
		if *req.Password != *req.ConfirmPassword {
			return nil, errConst.ErrPasswordDoesNotMatch
		}

		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		password = string(hashedPassword)
	}

	userResult, err = us.repository.UserRepo().Update(ctx, &dto.UpdateRequest{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    &password,
		PhoneNumber: req.PhoneNumber,
	}, uuid)

	if err != nil {
		return nil, err
	}

	data = dto.UserResponse{
		UUID:        userResult.UUID,
		Name:        userResult.Name,
		Username:    userResult.Username,
		PhoneNumber: userResult.PhoneNumber,
		Email:       userResult.Email,
	}

	return &data, nil

}

func (us *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		PhoneNumber: userLogin.PhoneNumber,
		Email:       userLogin.Email,
		Role:        userLogin.Role,
	}

	return &data, nil
}

func (us *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := us.repository.UserRepo().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	data := dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}
	return &data, nil
}
