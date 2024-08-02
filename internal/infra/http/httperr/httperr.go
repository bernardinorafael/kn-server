package httperr

import (
	"encoding/json"
	"net/http"
)

type HttpErr struct {
	Error   string `json:"httperr"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func BadRequestError(w http.ResponseWriter, message string) {
	code := http.StatusBadRequest

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}

func UnauthorizedError(w http.ResponseWriter, message string) {
	code := http.StatusUnauthorized

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}

func InternalServerError(w http.ResponseWriter, message string) {
	code := http.StatusInternalServerError

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}

func NotFoundError(w http.ResponseWriter, message string) {
	code := http.StatusNotFound

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}

func ConflictError(w http.ResponseWriter, message string) {
	code := http.StatusConflict

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}

func UnprocessableEntityError(w http.ResponseWriter, message string) {
	code := http.StatusUnprocessableEntity

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErr{
		Code:    code,
		Error:   http.StatusText(code),
		Message: message,
	})
}
