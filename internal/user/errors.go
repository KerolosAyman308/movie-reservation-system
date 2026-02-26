package user

import (
	"errors"
	"movie/system/pkg"
	"net/http"
)

var (
	ErrDuplicateEmail = errors.New("a user with that email already exists")
	ErrUserNotFound   = errors.New("user not found")

	ErrDuplicateEmailAPI = pkg.APIError[*string]{
		StatusCode: http.StatusConflict,
		Message:    "a user with that email already exists",
		Code:       "DUPLICATE_EMAIL",
	}
)
