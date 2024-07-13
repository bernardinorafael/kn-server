package error

import (
	"encoding/json"
	"net/http"
)

type RestError struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewBadRequestError(w http.ResponseWriter, message string) {
	statusCode := http.StatusBadRequest

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Code:    statusCode,
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func NewUnauthorizedError(w http.ResponseWriter, message string) {
	statusCode := http.StatusUnauthorized

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Code:    statusCode,
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func NewInternalServerError(w http.ResponseWriter, message string) {
	statusCode := http.StatusInternalServerError

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Code:    statusCode,
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func NewNotFoundError(w http.ResponseWriter, message string) {
	statusCode := http.StatusNotFound

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Code:    statusCode,
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func NewConflictError(w http.ResponseWriter, message string) {
	statusCode := http.StatusConflict

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Code:    statusCode,
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func NewUnprocessableEntityError(w http.ResponseWriter, message string) {
	statusCode := http.StatusUnprocessableEntity

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Code:    statusCode,
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}
