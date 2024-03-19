package httperr

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidationField struct {
	Field string `json:"field"`
	Msg   string `json:"message"`
}

type HttpErr struct {
	Err         string            `json:"error"`
	Code        int               `json:"code"`
	Message     string            `json:"message"`
	Validations []ValidationField `json:"validations,omitempty"`
}

func NewFieldsErrorValidation(c *gin.Context, m string, v []ValidationField) {
	c.JSON(http.StatusUnprocessableEntity, HttpErr{
		Err:         http.StatusText(http.StatusUnprocessableEntity),
		Code:        http.StatusUnprocessableEntity,
		Message:     m,
		Validations: v,
	})
}

func NewBadRequestError(c *gin.Context, m string) {
	c.JSON(http.StatusBadRequest, &HttpErr{
		Code:    http.StatusBadRequest,
		Err:     http.StatusText(http.StatusBadRequest),
		Message: m,
	})
}

func NewUnauthorizedError(c *gin.Context, m string) {
	c.JSON(http.StatusUnauthorized, &HttpErr{
		Code:    http.StatusUnauthorized,
		Err:     http.StatusText(http.StatusUnauthorized),
		Message: m,
	})
}

func NewInternalServerError(c *gin.Context, m string) {
	c.JSON(http.StatusInternalServerError, &HttpErr{
		Code:    http.StatusInternalServerError,
		Err:     http.StatusText(http.StatusInternalServerError),
		Message: m,
	})
}

func NewNotFoundError(c *gin.Context, m string) {
	c.JSON(http.StatusNotFound, &HttpErr{
		Code:    http.StatusNotFound,
		Err:     http.StatusText(http.StatusNotFound),
		Message: m,
	})
}

func NewForbiddenError(c *gin.Context, m string) {
	c.JSON(http.StatusForbidden, &HttpErr{
		Code:    http.StatusForbidden,
		Err:     http.StatusText(http.StatusForbidden),
		Message: m,
	})
}

func NewConflictError(c *gin.Context, m string) {
	c.JSON(http.StatusConflict, &HttpErr{
		Code:    http.StatusConflict,
		Err:     http.StatusText(http.StatusConflict),
		Message: m,
	})
}
