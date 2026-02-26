package user

import "context"

// Repository defines the data-access contract for the user domain.
// All database interactions must go through this interface.
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
}
