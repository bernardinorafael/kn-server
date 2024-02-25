package controller

import (
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/gozinho/internal/application/dto"
	"github.com/bernardinorafael/gozinho/internal/application/interfaces"
	"github.com/bernardinorafael/gozinho/internal/infra/http/error"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service interfaces.UserService
}

func NewUserController(r *gin.Engine, service interfaces.UserService) *UserController {
	c := &UserController{service}

	r.POST("/user", c.Save)
	r.GET("/user/:id", c.GetByID)
	r.PATCH("/user/:id", c.Update)
	r.DELETE("/user/:id", c.Delete)
	r.GET("/user", c.GetAll)
	r.PUT("/user/:id", c.UpdatePassword)

	return c
}

func (s UserController) Save(c *gin.Context) {
	req := dto.CreateUser{}

	if c.Request.Body == http.NoBody {
		slog.Error("body is required")
		httperror.NewBadRequestError(c, "body is required")
		return
	}

	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to decode body")
		httperror.NewBadRequestError(c, "failed to decode body")
		return
	}

	err = s.service.Save(&req)
	if err != nil {
		httperror.NewBadRequestError(c, "error creating user")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "ok"})
}

func (s UserController) GetByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		slog.Error("ID param not found")
		httperror.NewBadRequestError(c, "ID param not found")
		return
	}

	user, err := s.service.GetByID(id)
	if err != nil {
		slog.Error("error to get user", err)
		if err.Error() == "user not found" {
			httperror.NewNotFoundError(c, "user not found")
			return
		}
		httperror.NewBadRequestError(c, "error to get user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (s UserController) Update(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (s UserController) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		slog.Error("ID param not found")
		httperror.NewBadRequestError(c, "ID param not found")
		return
	}

	err := s.service.Delete(id)
	if err != nil {
		slog.Error("error to get user", err)
		if err.Error() == "user not found" {
			httperror.NewNotFoundError(c, "user not found")
			return
		}
		httperror.NewBadRequestError(c, "error to get user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

func (s UserController) GetAll(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (s UserController) UpdatePassword(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}
