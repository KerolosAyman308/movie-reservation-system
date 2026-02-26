package user

import (
	"context"

	dtos "movie/system/internal/user/DTOs"
)

// Service defines the business-logic contract for the user domain.
// All handlers must depend on this interface, never on a concrete type.
type Service interface {
	AddUser(ctx context.Context, dto dtos.UserCreateDTO) (*dtos.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*dtos.UserResponse, error)
	// GetUserByEmail returns the full User model (including password hash) for auth use.
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetAllUsers(ctx context.Context) ([]dtos.UserResponse, error)
	ChangeRole(ctx context.Context, id uint, dto dtos.ChangeRoleDTO) error
}
