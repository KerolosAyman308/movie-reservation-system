package user

import (
	"context"
	dtos "movie/system/internal/user/DTOs"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (s *UserService) AddUser(ctx context.Context, createDto dtos.UserCreateDTO) (*dtos.UserResponse, error) {
	user := User{
		Email:     createDto.Email,
		Password:  createDto.Password,
		FirstName: createDto.FirstName,
		LastName:  createDto.LastName,
		IsAdmin:   *createDto.IsAdmin,
		Birthday:  createDto.Birthday,
	}
	hashErr := user.SetPassword()
	if hashErr != nil {
		return nil, hashErr
	}

	userExist, _ := s.GetUserByEmail(ctx, createDto.Email)
	if userExist != nil {
		return nil, ErrDuplicateEmail
	}

	err := s.DB.Model(&user).WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &dtos.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) (*[]dtos.UserResponse, error) {
	var user []User
	err := s.DB.Model(&user).WithContext(ctx).Find(&user).Error

	if err != nil {
		return nil, err
	}

	// map to the DTO
	var dtoToReturn []dtos.UserResponse = make([]dtos.UserResponse, len(user))
	for i, value := range user {
		dtoToReturn[i] = dtos.UserResponse{
			ID:    value.ID,
			Email: value.Email,
		}
	}
	return &dtoToReturn, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*dtos.UserResponse, error) {
	var user User
	err := s.DB.Model(&user).WithContext(ctx).First(&user, id).Error

	if err != nil {
		return nil, userNotFound(id)
	}

	return &dtos.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*dtos.UserResponse, error) {
	var user User
	err := s.DB.Model(&user).WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &dtos.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (s *UserService) ChangeRole(ctx context.Context, userId uint, role dtos.ChangeRoleDTO) error {
	var user User
	err := s.DB.Model(&user).WithContext(ctx).Where("ID = ?", userId).First(&user).Error
	if err != nil {
		return userNotFound(userId)
	}

	user.IsAdmin = *role.IsAdmin
	updateErr := s.DB.Model(&user).WithContext(ctx).Where("ID = ?", userId).Save(user).Error

	return updateErr
}
