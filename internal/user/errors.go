package user

import (
	"errors"
	"fmt"
	"movie/system/pkg"
	"net/http"
)

var (
	ErrDuplicateEmailAPI pkg.APIError[*string] = pkg.APIError[*string]{StatusCode: http.StatusConflict, Message: ErrDuplicateEmail.Error(), Code: "ErrDuplicateEmail"}
)

var (
	ErrDuplicateEmail = errors.New("a user with that email already exists")
	ErrUserNotFound   = errors.New("user not found")
)

func userNotFound(id uint) error {
	return fmt.Errorf("could not find the requested user id %d: %w", id, ErrUserNotFound)
}

func ErrUserNotFoundAPI(err error) pkg.APIError[any] {
	return pkg.APIError[any]{
		StatusCode: http.StatusNotFound,
		Code:       "UserNotFound",
		Message:    err.Error(),
	}
}
