package errors

import "net/http"

type Fields struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

type HttpError struct {
	Message string   `json:"message"`
	Err     string   `json:"error,omitempty"`
	Code    int      `json:"code"`
	Fields  []Fields `json:"fields"`
}

func (h *HttpError) Error() string {
	return h.Message
}

func NewHttpError(m, e string, c int, f []Fields) *HttpError {
	return &HttpError{Message: m, Err: e, Code: c, Fields: f}
}

func NewBadRequestError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusBadRequest,
		Err:     "bad_request",
	}
}

func NewUnauthorizedRequestError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusUnauthorized,
		Err:     "unauthorized",
	}
}

func NewBadRequestValidationError(m string, f []Fields) *HttpError {
	return &HttpError{
		Message: m,
		Code:    http.StatusUnauthorized,
		Err:     "bad_request",
		Fields:  f,
	}
}

func NewInternalServerError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusInternalServerError,
		Err:     "internal_server_error",
	}
}

func NewNotFoundError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusNotFound,
		Err:     "not_found",
	}
}

func NewForbiddenError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusForbidden,
		Err:     "forbidden",
	}
}
