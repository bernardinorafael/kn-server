package restutil

import (
	"encoding/json"
	"net/http"
)

type ValidationField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type RestError struct {
	Error       string            `json:"error"`
	Code        int               `json:"code"`
	Message     string            `json:"message"`
	Validations []ValidationField `json:"validations,omitempty"`
}

func NewFieldsErrorValidation(w http.ResponseWriter, message string, validator []ValidationField) {
	statusCode := http.StatusUnprocessableEntity

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Error:       http.StatusText(statusCode),
		Code:        statusCode,
		Message:     message,
		Validations: validator,
	})
}

func NewBadRequestError(w http.ResponseWriter, message string) {
	statusCode := http.StatusBadRequest

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	})
}

func NewUnauthorizedError(w http.ResponseWriter, message string) {
	statusCode := http.StatusUnauthorized

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	})
}

func NewInternalServerError(w http.ResponseWriter, message string) {
	statusCode := http.StatusInternalServerError

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	})
}

func NewNotFoundError(w http.ResponseWriter, message string) {
	statusCode := http.StatusNotFound

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	})
}

func NewConflictError(w http.ResponseWriter, message string) {
	statusCode := http.StatusConflict

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	})
}

func NewUnprocessableEntityError(w http.ResponseWriter, message string) {
	statusCode := http.StatusUnprocessableEntity

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(RestError{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	})
}
