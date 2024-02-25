package controllers

import (
	"github.com/bernardinorafael/gozinho/internal/application/interfaces"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service interfaces.UserService
}

func NewUserController(r *gin.Engine, service interfaces.UserService) *UserController {
	controller := &UserController{service}

	r.POST("/user", controller.CreateUser)

	return controller
}

func (s *UserController) CreateUser(c *gin.Context) {
	// handle controller here
}
