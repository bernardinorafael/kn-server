package accountroute

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/bernardinorafael/gozinho/internal/application/contract"
	"github.com/bernardinorafael/gozinho/internal/application/dto"
	resterror "github.com/bernardinorafael/gozinho/internal/infra/rest/error"
	"github.com/bernardinorafael/gozinho/internal/infra/rest/restutil"
	"github.com/gin-gonic/gin"
)

var handler *Handler
var once sync.Once

type Handler struct {
	svc contract.AccountService
}

func NewHandler(svc contract.AccountService) *Handler {
	once.Do(func() {
		handler = &Handler{
			svc: svc,
		}
	})

	return handler
}

func (s Handler) Save(c *gin.Context) {
	ctx := restutil.GetContext(c)

	input := dto.UserInput{}
	if c.Request.Body == http.NoBody {
		resterror.NewBadRequestError(c, "body is required")
		return
	}

	err := c.ShouldBind(&input)
	if err != nil {
		resterror.NewBadRequestError(c, "failed to decode body")
		return
	}

	err = s.svc.Save(ctx, &input)
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
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	user, err := s.svc.GetByID(ctx, id)
	if err != nil {
		resterror.NewNotFoundError(c, "user not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (s Handler) Update(c *gin.Context) {
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	input := dto.UpdateUser{}

	if c.Request.Body == http.NoBody {
		resterror.NewBadRequestError(c, "body is required")
		return
	}

	err := c.ShouldBind(&input)
	if err != nil {
		resterror.NewBadRequestError(c, "failed to decode body")
		return
	}

	err = s.svc.Update(ctx, &input, id)
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
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	if id == "" {
		resterror.NewBadRequestError(c, "ID param not found")
		return
	}

	err := s.svc.Delete(ctx, id)
	if err != nil {
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
	ctx := restutil.GetContext(c)

	users, err := s.svc.GetAll(ctx)
	if err != nil {
		resterror.NewBadRequestError(c, "error to get all users")
		return
	}

	c.JSON(http.StatusOK, users)
}

func (s Handler) UpdatePassword(c *gin.Context) {
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	input := dto.UpdatePassword{}

	if c.Request.Body == http.NoBody {
		resterror.NewBadRequestError(c, "body is required")
		return
	}

	err := c.ShouldBind(&input)
	if err != nil {
		resterror.NewBadRequestError(c, "failed to decode body")
		return
	}

	err = s.svc.UpdatePassword(ctx, &input, id)
	if err != nil {
		resterror.NewConflictError(c, "error on update password")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
