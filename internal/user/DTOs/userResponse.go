package dtos

import "time"

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  *string   `json:"last_name,omitempty"`
	IsAdmin   bool      `json:"is_admin"`
	Birthday  time.Time `json:"birthday"`
}
