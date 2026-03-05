package user

import (
	"context"
	"testing"
	"time"

	dtos "movie/system/internal/user/DTOs"
)

type mockRepo struct {
	users []User
}

func (m *mockRepo) Create(ctx context.Context, u *User) error {
	u.ID = uint(len(m.users) + 1)
	m.users = append(m.users, *u)
	return nil
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *mockRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *mockRepo) FindAll(ctx context.Context) ([]User, error) {
	return m.users, nil
}

func (m *mockRepo) Update(ctx context.Context, user *User) error {
	for i, u := range m.users {
		if u.ID == user.ID {
			m.users[i] = *user
			return nil
		}
	}
	return ErrUserNotFound
}

func TestUserService_AddUser(t *testing.T) {
	repo := &mockRepo{}
	service := NewService(repo)

	isAdmin := false
	dto := dtos.UserCreateDTO{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  func(s string) *string { return &s }("Doe"),
		IsAdmin:   &isAdmin,
		Birthday:  time.Now(),
	}

	res, err := service.AddUser(context.Background(), dto)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}
	if res.Email != dto.Email {
		t.Errorf("expected email %s, got %s", dto.Email, res.Email)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	repo := &mockRepo{}
	service := NewService(repo)

	// Pre-seed the repository
	repo.users = append(repo.users, User{
		ID:        1,
		Email:     "existing@example.com",
		FirstName: "Jane",
	})

	res, err := service.GetUserByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if res.Email != "existing@example.com" {
		t.Errorf("Expected email existing@example.com, got %s", res.Email)
	}
}
