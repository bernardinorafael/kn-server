package authroute

import (
	"errors"
	"net/http"
	"sync"

	httperr "github.com/bernardinorafael/kn-server/helper/error"
	"github.com/bernardinorafael/kn-server/helper/validator"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/application/service"
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
	req := dto.Login{}

	err := c.ShouldBind(&req)
	if err != nil {
		httperr.NewBadRequestError(c, "not found/invalid body")
		return
	}

	val := validator.Validate(req)
	if val != nil {
		httperr.NewFieldsErrorValidation(c, "invalid fields", val)
		return
	}

	user, err := h.authService.Login(req)
	if err != nil {
		httperr.NewUnauthorizedError(c, err.Error())
		return
	}

	token, claims, err := h.jwtService.CreateToken(user.ID)
	if err != nil {
		httperr.NewUnauthorizedError(c, err.Error())
		return
	}

	r := LoginResponse{
		UserID:      claims.UserID,
		IssuedAt:    claims.IssuedAt,
		ExpiresAt:   claims.ExpiresAt,
		AccessToken: token,
	}

	c.JSON(http.StatusOK, r)
}

func (h *UserHandler) Register(c *gin.Context) {
	req := dto.Register{}

	err := c.ShouldBind(&req)
	if err != nil {
		httperr.NewBadRequestError(c, "not found/invalid body")
		return
	}

	validations := validator.Validate(req)
	if validations != nil {
		httperr.NewFieldsErrorValidation(c, "invalid fields", validations)
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyTaken) {
			httperr.NewConflictError(c, err.Error())
			return
		} else if errors.Is(err, service.ErrDocumentAlreadyTaken) {
			httperr.NewConflictError(c, err.Error())
			return
		}
		httperr.NewBadRequestError(c, err.Error())
		return
	}

	token, claims, err := h.jwtService.CreateToken(user.ID)
	if err != nil {
		httperr.NewBadRequestError(c, err.Error())
		return
	}

	r := LoginResponse{
		UserID:      claims.UserID,
		IssuedAt:    claims.IssuedAt,
		ExpiresAt:   claims.ExpiresAt,
		AccessToken: token,
	}

	c.JSON(http.StatusCreated, r)
}
