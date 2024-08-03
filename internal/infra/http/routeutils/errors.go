package routeutils

import (
	"encoding/json"
	"net/http"
)

type HttpErr struct {
	Error   string `json:"httperr"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewBadRequestError(w http.ResponseWriter, message string) {
	code := http.StatusBadRequest

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}

func NewUnauthorizedError(w http.ResponseWriter, message string) {
	code := http.StatusUnauthorized

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}

func NewInternalServerError(w http.ResponseWriter, message string) {
	code := http.StatusInternalServerError

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}

func NewNotFoundError(w http.ResponseWriter, message string) {
	code := http.StatusNotFound

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}

func NewConflictError(w http.ResponseWriter, message string) {
	code := http.StatusConflict

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}

func NewUnprocessableEntityError(w http.ResponseWriter, message string) {
	code := http.StatusUnprocessableEntity

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}
