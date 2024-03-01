package accountroute

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	resterror "github.com/bernardinorafael/kn-server/internal/infra/rest/error"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/response"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
	"github.com/gin-gonic/gin"
)

var (
	handler *Handler
	once    sync.Once
)

type Handler struct {
	as   contract.AccountService
	auth contract.AuthService
}

func New(as contract.AccountService, auth contract.AuthService) *Handler {
	once.Do(func() {
		handler = &Handler{as, auth}
	})
	return handler
}

func (h Handler) Login(c *gin.Context) {
	ctx := restutil.GetContext(c)

	credentials := dto.Login{}
	if c.Request.Body == http.NoBody {
		resterror.NewBadRequestError(c, "not found/invalid body")
		return
	}

	err := c.ShouldBind(&credentials)
	if err != nil {
		resterror.NewBadRequestError(c, "not found/invalid body")
		return
	}

	account, err := h.as.Login(ctx, credentials)
	if err != nil {
		resterror.NewUnauthorizedRequestError(c, "failed to login")
		return
	}

	input := dto.TokenPayloadInput{ID: account.ID}
	token, payload, err := h.auth.CreateAccessToken(ctx, input)
	if err != nil {
		resterror.NewBadRequestError(c, err.Error())
	}

	r := response.LoginResponse{
		AccessToken: token,
		UserID:      payload.UserID,
		ExpiresAt:   payload.ExpiresAt,
		IssuedAt:    payload.ExpiresAt,
	}

	c.JSON(http.StatusOK, r)
}

func (h Handler) Save(c *gin.Context) {
	ctx := restutil.GetContext(c)

	input := dto.UserInput{}
	if c.Request.Body == http.NoBody {
		resterror.NewBadRequestError(c, "not found/invalid body")
		return
	}

	err := c.ShouldBind(input)
	if err != nil {
		resterror.NewBadRequestError(c, "not found/invalid body")
		return
	}

	err = h.as.Save(ctx, input)
	if err != nil {
		if err.Error() == "user already taken" {
			resterror.NewConflictError(c, "credential(s) already taken")
			return
		}
		resterror.NewBadRequestError(c, "error creating user")
		return
	}
	c.Status(http.StatusCreated)
}

func (h Handler) GetByID(c *gin.Context) {
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	user, err := h.as.GetByID(ctx, id)
	if err != nil {
		resterror.NewNotFoundError(c, "user not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h Handler) Update(c *gin.Context) {
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	input := dto.UpdateUser{}

	if c.Request.Body == http.NoBody {
		resterror.NewBadRequestError(c, "not found/invalid body")
		return
	}

	err := c.ShouldBind(input)
	if err != nil {
		resterror.NewBadRequestError(c, "not found/invalid body")
		return
	}

	err = h.as.Update(ctx, input, id)
	if err != nil {
		slog.Error("error to update user", err)
		if err.Error() == "user not found" {
			resterror.NewNotFoundError(c, "user not found")
			return
		}
		if err.Error() == "user already exists" {
			resterror.NewConflictError(c, "email already taken")
		}
		resterror.NewBadRequestError(c, "error to get user")
		return
	}
	c.Status(http.StatusOK)
}

func (h Handler) Delete(c *gin.Context) {
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	if id == "" {
		resterror.NewBadRequestError(c, "ID param not found")
		return
	}

	err := h.as.Delete(ctx, id)
	if err != nil {
		if err.Error() == "user not found" {
			resterror.NewNotFoundError(c, "user not found")
			return
		}
		resterror.NewBadRequestError(c, "error to get user")
		return
	}
	c.Status(http.StatusOK)
}

func (h Handler) GetAll(c *gin.Context) {
	ctx := restutil.GetContext(c)

	users, err := h.as.GetAll(ctx)
	if err != nil {
		resterror.NewBadRequestError(c, "error to get all users")
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h Handler) UpdatePassword(c *gin.Context) {
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	input := dto.UpdatePassword{}

	if c.Request.Body == http.NoBody {
		resterror.NewBadRequestError(c, "not found/invalid body")
		return
	}

	err := c.ShouldBind(input)
	if err != nil {
		resterror.NewBadRequestError(c, "not found/invalid body")
		return
	}

	err = h.as.UpdatePassword(ctx, input, id)
	if err != nil {
		resterror.NewConflictError(c, "error on update password")
		return
	}
	c.Status(http.StatusOK)
}
