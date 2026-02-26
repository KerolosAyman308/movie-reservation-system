package pkg

import (
	log "log/slog"
	"net/http"
)

type APIResponse[T any] struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

type APIError[T any] struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Data       *T     `json:"data"`
}

func Ok[T any](data T, message string, w http.ResponseWriter) {
	if message == "" {
		message = "Success"
	}

	resp := APIResponse[T]{
		StatusCode: http.StatusOK,
		Data:       data,
		Message:    message,
	}

	WriteJSON(w, 200, resp)
}

func Error[T any](err APIError[T], w http.ResponseWriter, r *http.Request) {

	log.Error(err.Message, "Code", err.Code, "status", err.StatusCode, "data", err.Data, "method", r.Method, "path", r.URL.Path)
	WriteJSON(w, err.StatusCode, err)
}
