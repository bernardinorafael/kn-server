package resterr

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Fields struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

type HttpResponseErr struct {
	Message string   `json:"message"`
	Err     string   `json:"error"`
	Code    int      `json:"code"`
	Fields  []Fields `json:"fields,omitempty"`
}

func NewBadRequestError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, HttpResponseErr{
		Message: msg,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
	})
}

func NewUnauthorizedError(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, HttpResponseErr{
		Message: msg,
		Err:     "unauthorized",
		Code:    http.StatusUnauthorized,
	})
}

func NewBadRequestValidationError(c *gin.Context, msg string, fields []Fields) {
	c.JSON(http.StatusBadRequest, HttpResponseErr{
		Message: msg,
		Err:     "bad_request_validation",
		Code:    http.StatusBadRequest,
		Fields:  fields,
	})
}

func NewInternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, HttpResponseErr{
		Message: msg,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	})
}

func NewNotFoundError(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, HttpResponseErr{
		Message: msg,
		Err:     "not_found",
		Code:    http.StatusNotFound,
	})
}

func NewForbiddenError(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, HttpResponseErr{
		Message: msg,
		Err:     "forbidden",
		Code:    http.StatusForbidden,
	})
}

func NewConflictError(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, HttpResponseErr{
		Message: msg,
		Err:     "conflict",
		Code:    http.StatusConflict,
	})
}
