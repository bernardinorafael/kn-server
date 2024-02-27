package accountroute

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/bernardinorafael/gozinho/internal/application/contract"
	"github.com/bernardinorafael/gozinho/internal/application/dto"
	resterror "github.com/bernardinorafael/gozinho/internal/infra/rest/error"
	"github.com/gin-gonic/gin"
)

var handler *Handler
var once sync.Once

type Handler struct {
	accountService contract.AccountService
}

func NewHandler(accountService contract.AccountService) *Handler {
	once.Do(func() {
		handler = &Handler{
			accountService: accountService,
		}
	})

	return handler
}

func (s Handler) Save(c *gin.Context) {
	input := dto.UserInput{}

	if c.Request.Body == http.NoBody {
		slog.Error("body is required")
		resterror.NewBadRequestError(c, "body is required")
		return
	}

	err := c.ShouldBind(&input)
	if err != nil {
		slog.Error("failed to decode body")
		resterror.NewBadRequestError(c, "failed to decode body")
		return
	}

	err = s.accountService.Save(&input)
	if err != nil {
		if err.Error() == "email already taken" {
			resterror.NewConflictError(c, err.Error())
			return
		}
		resterror.NewBadRequestError(c, "error creating user")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

func (s Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := s.accountService.GetByID(id)
	if err != nil {
		slog.Error("error to get user", err, slog.String("pkg", "controller"))
		resterror.NewNotFoundError(c, "user not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (s Handler) Update(c *gin.Context) {
	id := c.Param("id")
	input := dto.UpdateUser{}

	if c.Request.Body == http.NoBody {
		slog.Error("body is required")
		resterror.NewBadRequestError(c, "body is required")
		return
	}

	err := c.ShouldBind(&input)
	if err != nil {
		slog.Error("failed to decode body")
		resterror.NewBadRequestError(c, "failed to decode body")
		return
	}

	err = s.accountService.Update(&input, id)
	if err != nil {
		slog.Error("error to update user", err)
		if err.Error() == "user not found" {
			resterror.NewNotFoundError(c, "user not found")
			return
		}
		if err.Error() == "user already exists" {
			resterror.NewConflictError(c, "e-mail already taken")
		}
		resterror.NewBadRequestError(c, "error to get user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		slog.Error("ID param not found")
		resterror.NewBadRequestError(c, "ID param not found")
		return
	}

	err := s.accountService.Delete(id)
	if err != nil {
		slog.Error("error to get user", err)
		if err.Error() == "user not found" {
			resterror.NewNotFoundError(c, "user not found")
			return
		}
		resterror.NewBadRequestError(c, "error to get user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s Handler) GetAll(c *gin.Context) {
	users, err := s.accountService.GetAll()
	if err != nil {
		slog.Error("error to get all users", err)
		resterror.NewBadRequestError(c, "error to get all users")
		return
	}

	c.JSON(http.StatusOK, users)
}

func (s Handler) UpdatePassword(c *gin.Context) {
	id := c.Param("id")
	input := dto.UpdatePassword{}

	if c.Request.Body == http.NoBody {
		slog.Error("body is required", slog.String("pkg", "controller"))
		resterror.NewBadRequestError(c, "body is required")
		return
	}

	err := c.ShouldBind(&input)
	if err != nil {
		slog.Error("failed to decode body", slog.String("pkg", "controller"))
		resterror.NewBadRequestError(c, "failed to decode body")
		return
	}

	err = s.accountService.UpdatePassword(&input, id)
	if err != nil {
		slog.Error("error to update password", err, slog.String("pkg", "controller"))
		resterror.NewConflictError(c, "error on update password")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
