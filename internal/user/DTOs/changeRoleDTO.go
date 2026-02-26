package dtos

type ChangeRoleDTO struct {
	IsAdmin *bool `json:"is_admin" validate:"required"`
}
