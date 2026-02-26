package auth

import (
	"movie/system/pkg"
	"net/http"
)

var (
	ErrNotAuthorized = pkg.APIError[string]{StatusCode: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Email or password is wrong"}
)
