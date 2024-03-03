package httperr

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type httpErr struct {
	Err     string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHttpErr(err, msg string, code int) httpErr {
	return httpErr{
		Code:    code,
		Err:     err,
		Message: msg,
	}
}

func NewBadRequestError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, httpErr{
		Code:    http.StatusBadRequest,
		Err:     http.StatusText(http.StatusBadRequest),
		Message: msg,
	})
}

func NewUnauthorizedError(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, httpErr{
		Code:    http.StatusUnauthorized,
		Err:     http.StatusText(http.StatusUnauthorized),
		Message: msg,
	})
}

func NewInternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, httpErr{
		Code:    http.StatusInternalServerError,
		Err:     http.StatusText(http.StatusInternalServerError),
		Message: msg,
	})
}

func NewNotFoundError(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, httpErr{
		Code:    http.StatusNotFound,
		Err:     http.StatusText(http.StatusNotFound),
		Message: msg,
	})
}

func NewForbiddenError(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, httpErr{
		Code:    http.StatusForbidden,
		Err:     http.StatusText(http.StatusForbidden),
		Message: msg,
	})
}

func NewConflictError(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, httpErr{
		Code:    http.StatusConflict,
		Err:     http.StatusText(http.StatusConflict),
		Message: msg,
	})
}
