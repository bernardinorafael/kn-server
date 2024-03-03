package authroute

import (
	"errors"
	"net/http"
	"sync"

	httperr "github.com/bernardinorafael/kn-server/helper/error"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
	"github.com/gin-gonic/gin"
)

var handler *UserHandler
var once sync.Once

type UserHandler struct {
	authService contract.AuthService
	jwtService  contract.JWTService
}

func NewUserHandler(as contract.AuthService, js contract.JWTService) *UserHandler {
	once.Do(func() {
		handler = &UserHandler{
			jwtService:  js,
			authService: as,
		}
	})
	return handler
}

func (h *UserHandler) Login(c *gin.Context) {
	ctx := restutil.GetContext(c)

	req := dto.Login{}
	err := c.ShouldBind(&req)
	if err != nil {
		httperr.NewBadRequestError(c, "not found/invalid body")
		return
	}

	id, err := h.authService.Login(ctx, req)
	if err != nil {
		httperr.NewUnauthorizedError(c, err.Error())
		return
	}

	token, claims, err := h.jwtService.CreateToken(ctx, id)
	if err != nil {
		httperr.NewUnauthorizedError(c, err.Error())
		return
	}

	r := LoginResponse{
		AccessToken: token,
		UserID:      claims.UserID,
		ExpiresAt:   claims.ExpiresAt,
		IssuedAt:    claims.IssuedAt,
	}

	c.JSON(http.StatusOK, r)
}

func (h *UserHandler) Register(c *gin.Context) {
	ctx := restutil.GetContext(c)

	req := dto.Register{}
	err := c.ShouldBind(&req)
	if err != nil {
		httperr.NewBadRequestError(c, "not found/invalid body")
		return
	}

	id, err := h.authService.Register(ctx, req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyTaken) {
			httperr.NewConflictError(c, err.Error())
			return
		}
		httperr.NewBadRequestError(c, err.Error())
		return
	}

	token, claims, err := h.jwtService.CreateToken(ctx, id)
	if err != nil {
		httperr.NewBadRequestError(c, err.Error())
		return
	}

	r := LoginResponse{
		AccessToken: token,
		UserID:      claims.UserID,
		ExpiresAt:   claims.ExpiresAt,
		IssuedAt:    claims.IssuedAt,
	}

	c.JSON(http.StatusCreated, r)
}
