package accountroute

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	resterr "github.com/bernardinorafael/kn-server/internal/infra/rest/error"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/response"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
	"github.com/gin-gonic/gin"
)

var (
	handler *Handler
	once    sync.Once
)

type Handler struct {
	service contract.AccountService
	auth    contract.AuthService
}

func New(service contract.AccountService, auth contract.AuthService) *Handler {
	once.Do(func() {
		handler = &Handler{service, auth}
	})
	return handler
}

func (h Handler) Login(c *gin.Context) {
	ctx := restutil.GetContext(c)

	credentials := dto.Login{}
	if c.Request.Body == http.NoBody {
		resterr.NewBadRequestError(c, "not found/invalid body")
		return
	}

	err := c.ShouldBind(&credentials)
	if err != nil {
		resterr.NewBadRequestError(c, "not found/invalid body")
		return
	}

	account, err := h.service.Login(ctx, credentials)
	if err != nil {
		resterr.NewUnauthorizedError(c, "failed to login")
		return
	}

	userID := dto.TokenPayloadInput{ID: account.ID}
	token, payload, err := h.auth.CreateAccessToken(ctx, userID)
	if err != nil {
		resterr.NewBadRequestError(c, err.Error())
	}

	r := response.LoginResponse{
		AccessToken: token,
		UserID:      payload.UserID,
		ExpiresAt:   payload.ExpiresAt,
		IssuedAt:    payload.IssuedAt,
	}

	c.JSON(http.StatusOK, r)
}

func (h Handler) Save(c *gin.Context) {
	ctx := restutil.GetContext(c)

	input := dto.UserInput{}
	err := c.ShouldBind(&input)
	if err != nil {
		resterr.NewBadRequestError(c, "not found/invalid body")
		return
	}

	userID, err := h.service.Save(ctx, input)
	if err != nil {
		resterr.NewBadRequestError(c, "error creating user")
		return
	}

	id := dto.TokenPayloadInput{ID: userID}
	token, payload, err := h.auth.CreateAccessToken(ctx, id)
	if err != nil {
		resterr.NewBadRequestError(c, err.Error())
	}

	r := response.LoginResponse{
		AccessToken: token,
		UserID:      payload.UserID,
		ExpiresAt:   payload.ExpiresAt,
		IssuedAt:    payload.IssuedAt,
	}

	c.JSON(http.StatusOK, r)

}

func (h Handler) GetByID(c *gin.Context) {
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	user, err := h.service.GetByID(ctx, id)
	if err != nil {
		resterr.NewNotFoundError(c, "user not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h Handler) Update(c *gin.Context) {
	ctx := restutil.GetContext(c)
	id := c.Param("id")

	input := dto.UpdateUser{}
	err := c.ShouldBind(input)
	if err != nil {
		resterr.NewBadRequestError(c, "invalid body request")
		return
	}

	err = h.service.Update(ctx, input, id)
	if err != nil {
		slog.Error("error to update user", err)
		if err.Error() == "user not found" {
			resterr.NewNotFoundError(c, "user not found")
			return
		}
		if err.Error() == "user already exists" {
			resterr.NewConflictError(c, "email already taken")
		}
		resterr.NewBadRequestError(c, "error to get user")
		return
	}

	c.Status(http.StatusOK)
}

func (h Handler) Delete(c *gin.Context) {
	ctx := restutil.GetContext(c)
	id := c.Param("id")

	err := h.service.Delete(ctx, id)
	if err != nil {
		resterr.NewNotFoundError(c, "the user you are trying to delete was not found")
		return
	}

	c.Status(http.StatusOK)
}

func (h Handler) GetAll(c *gin.Context) {
	ctx := restutil.GetContext(c)

	accounts, err := h.service.GetAll(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "the list is empty") {
			c.JSON(http.StatusOK, gin.H{
				"message": "there are no users registered at the moment",
			})
			return
		}
		resterr.NewBadRequestError(c, "error to get all users")
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func (h Handler) UpdatePassword(c *gin.Context) {
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	input := dto.UpdatePassword{}

	if c.Request.Body == http.NoBody {
		resterr.NewBadRequestError(c, "not found/invalid body")
		return
	}

	err := c.ShouldBind(input)
	if err != nil {
		resterr.NewBadRequestError(c, "not found/invalid body")
		return
	}

	err = h.service.UpdatePassword(ctx, input, id)
	if err != nil {
		resterr.NewConflictError(c, "error on update password")
		return
	}
	c.Status(http.StatusOK)
}
