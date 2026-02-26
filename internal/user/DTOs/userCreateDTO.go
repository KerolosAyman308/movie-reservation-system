package dtos

import "time"

type UserCreateDTO struct {
	Email     string    `json:"email" validate:"required,email"`
	FirstName string    `json:"first_name" validate:"required,min=3,max=50"`
	LastName  *string   `json:"last_name,omitempty" validate:"omitempty,max=50"`
	IsAdmin   *bool     `json:"is_admin" validate:"required"`
	Birthday  time.Time `json:"birthday" validate:"required"`
	Password  string    `json:"Password" validate:"required,min=8,max=20"`
}

type UserSignUpDTO struct {
	Email     string    `json:"email" validate:"required,email"`
	FirstName string    `json:"first_name" validate:"required,min=3,max=50"`
	LastName  *string   `json:"last_name,omitempty" validate:"omitempty,max=50"`
	Birthday  time.Time `json:"birthday" validate:"required,datetime"`
	Password  string    `json:"Password" validate:"required,min=8,max=20"`
}

type UserLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"Password" validate:"required,min=8,max=20"`
}
