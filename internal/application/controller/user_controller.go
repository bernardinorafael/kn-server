package controller

import (
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/gozinho/internal/application/contract"
	"github.com/bernardinorafael/gozinho/internal/application/dto"
	"github.com/bernardinorafael/gozinho/internal/infra/http/error"
	"github.com/gin-gonic/gin"
)

type userController struct {
	svc contract.UserService
}

func NewUserController(r *gin.Engine, svc contract.UserService) *userController {
	c := &userController{svc}

	r.POST("/user", c.Save)
	r.GET("/user/:id", c.GetByID)
	r.PATCH("/user/:id", c.Update)
	r.DELETE("/user/:id", c.Delete)
	r.GET("/users", c.GetAll)
	r.PUT("/user/:id", c.UpdatePassword)

	return c
}

func (s userController) Save(c *gin.Context) {
	req := dto.UserInput{}

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

	err = s.svc.Save(&req)
	if err != nil {
		if err.Error() == "email already taken" {
			httperror.NewConflictError(c, err.Error())
			return
		}
		httperror.NewBadRequestError(c, "error creating user")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

func (s userController) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := s.svc.GetByID(id)
	if err != nil {
		slog.Error("error to get user", err, slog.String("pkg", "controller"))
		httperror.NewNotFoundError(c, "user not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (s userController) Update(c *gin.Context) {
	id := c.Param("id")
	req := dto.UpdateUser{}

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

	err = s.svc.Update(&req, id)
	if err != nil {
		slog.Error("error to update user", err)
		if err.Error() == "user not found" {
			httperror.NewNotFoundError(c, "user not found")
			return
		}
		if err.Error() == "user already exists" {
			httperror.NewConflictError(c, "e-mail already taken")
		}
		httperror.NewBadRequestError(c, "error to get user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s userController) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		slog.Error("ID param not found")
		httperror.NewBadRequestError(c, "ID param not found")
		return
	}

	err := s.svc.Delete(id)
	if err != nil {
		slog.Error("error to get user", err)
		if err.Error() == "user not found" {
			httperror.NewNotFoundError(c, "user not found")
			return
		}
		httperror.NewBadRequestError(c, "error to get user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s userController) GetAll(c *gin.Context) {
	users, err := s.svc.GetAll()
	if err != nil {
		slog.Error("error to get all users", err)
		httperror.NewBadRequestError(c, "error to get all users")
		return
	}

	c.JSON(http.StatusOK, users)
}

func (s userController) UpdatePassword(c *gin.Context) {
	id := c.Param("id")
	req := dto.UpdatePassword{}

	if c.Request.Body == http.NoBody {
		slog.Error("body is required", slog.String("pkg", "controller"))
		httperror.NewBadRequestError(c, "body is required")
		return
	}

	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to decode body", slog.String("pkg", "controller"))
		httperror.NewBadRequestError(c, "failed to decode body")
		return
	}

	err = s.svc.UpdatePassword(&req, id)
	if err != nil {
		slog.Error("error to update password", err, slog.String("pkg", "controller"))
		httperror.NewConflictError(c, "error on update password")
		return

	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
