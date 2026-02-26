package pkg

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

func NotFound[T any](w http.ResponseWriter, r *http.Request, data *T) {
	err := APIError[T]{StatusCode: http.StatusNotFound, Message: "Resource not found", Code: "NOT_FOUND", Data: data}
	Error(err, w, r)
}

func Unauthorized[T any](w http.ResponseWriter, r *http.Request, data *T) {
	err := APIError[T]{StatusCode: http.StatusUnauthorized, Message: "Unauthorized access", Code: "UNAUTHORIZED", Data: data}
	Error(err, w, r)
}

func InternalError[T any](w http.ResponseWriter, r *http.Request, data *T) {
	err := APIError[T]{StatusCode: http.StatusInternalServerError, Message: "Internal server error", Code: "INTERNAL_ERROR", Data: data}
	Error(err, w, r)
}

func BadRequest[T any](w http.ResponseWriter, r *http.Request, data *T) {
	err := APIError[T]{StatusCode: http.StatusBadRequest, Message: "Bad Request", Code: "BAD_REQUEST", Data: data}
	Error(err, w, r)
}

func BadRequestWithCustomMessage[T any](w http.ResponseWriter, r *http.Request, message string, data *T) {
	err := APIError[T]{StatusCode: http.StatusBadRequest, Message: message, Code: "BAD_REQUEST", Data: data}
	Error(err, w, r)
}

func Conflict[T any](w http.ResponseWriter, r *http.Request, data *T) {
	err := APIError[T]{StatusCode: http.StatusConflict, Message: "Conflict Response", Code: "CONFLICT_RESPONSE", Data: data}
	Error(err, w, r)
}

func Forbidden[T any](w http.ResponseWriter, r *http.Request, data *T) {
	err := APIError[T]{StatusCode: http.StatusForbidden, Message: "Access denied", Code: "FORBIDDEN", Data: data}
	Error(err, w, r)
}

func ValidateStruct(s interface{}) map[string]string {
	if err := Validate.Struct(s); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMap := make(map[string]string)
		for _, err := range validationErrors {
			errorMap[err.Field()] = err.Error()
		}

		return errorMap
	}
	return nil
}
