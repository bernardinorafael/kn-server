package interfaces

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Save(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetAll(c *gin.Context)
	UpdatePassword(c *gin.Context)
}