package user

import (
	"context"

	dtos "movie/system/internal/user/DTOs"
)

// UserService is the concrete implementation of the Service interface.
type UserService struct {
	repo Repository
}

// NewService creates a UserService backed by the provided Repository.
func NewService(repo Repository) Service {
	return &UserService{repo: repo}
}

func (s *UserService) AddUser(ctx context.Context, dto dtos.UserCreateDTO) (*dtos.UserResponse, error) {
	u := User{
		Email:     dto.Email,
		Password:  dto.Password,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		IsAdmin:   *dto.IsAdmin,
		Birthday:  dto.Birthday,
	}
	if err := u.SetPassword(); err != nil {
		return nil, err
	}
	// Duplicate-email detection is handled at the DB layer via unique index —
	// no check-then-act race condition.
	if err := s.repo.Create(ctx, &u); err != nil {
		return nil, err
	}
	return &dtos.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsAdmin:   u.IsAdmin,
		Birthday:  u.Birthday,
	}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]dtos.UserResponse, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dtos.UserResponse, len(users))
	for i, u := range users {
		result[i] = dtos.UserResponse{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			IsAdmin:   u.IsAdmin,
			Birthday:  u.Birthday,
		}
	}
	return result, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*dtos.UserResponse, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dtos.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsAdmin:   u.IsAdmin,
		Birthday:  u.Birthday,
	}, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) ChangeRole(ctx context.Context, id uint, dto dtos.ChangeRoleDTO) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	u.IsAdmin = *dto.IsAdmin
	return s.repo.Update(ctx, u)
}
