package resterror

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Fields struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

type HttpResponseError struct {
	Message string   `json:"message"`
	Err     string   `json:"error"`
	Code    int      `json:"code"`
	Fields  []Fields `json:"fields,omitempty"`
}

func NewBadRequestError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, HttpResponseError{
		Message: msg,
		Err:     "BAD_REQUEST",
		Code:    http.StatusBadRequest,
	})
}

func NewUnauthorizedRequestError(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, HttpResponseError{
		Message: msg,
		Err:     "UNAUTHORIZED",
		Code:    http.StatusUnauthorized,
	})
}

func NewBadRequestValidationError(c *gin.Context, msg string, fields []Fields) {
	c.JSON(http.StatusBadRequest, HttpResponseError{
		Message: msg,
		Err:     "bad_request_validation",
		Code:    http.StatusBadRequest,
		Fields:  fields,
	})
}

func NewInternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, HttpResponseError{
		Message: msg,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	})
}

func NewNotFoundError(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, HttpResponseError{
		Message: msg,
		Err:     "not_found",
		Code:    http.StatusNotFound,
	})
}

func NewForbiddenError(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, HttpResponseError{
		Message: msg,
		Err:     "forbidden",
		Code:    http.StatusForbidden,
	})
}

func NewConflictError(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, HttpResponseError{
		Message: msg,
		Err:     "conflict",
		Code:    http.StatusConflict,
	})
}
