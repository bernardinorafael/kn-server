package route

import (
	"errors"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/http/middleware"
	"github.com/bernardinorafael/kn-server/internal/infra/http/response"
	. "github.com/bernardinorafael/kn-server/internal/infra/http/routeutils"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type userHandler struct {
	log         logger.Logger
	userService contract.UserService
	jwtAuth     auth.TokenAuthInterface
}

func NewUserHandler(log logger.Logger, userService contract.UserService, jwtAuth auth.TokenAuthInterface) *userHandler {
	return &userHandler{log, userService, jwtAuth}
}

func (h userHandler) RegisterRoute(r *chi.Mux) {
	m := middleware.NewWithAuth(h.jwtAuth, h.log)

	r.Route("/users", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Get("/me", h.getSigned)
		r.Put("/{id}", h.updateUser)
		r.Patch("/{id}/password", h.recoverPassword)
	})
}

func (h userHandler) recoverPassword(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdatePassword

	err := ParseBodyRequest(r, &payload)
	if err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	err = h.userService.RecoverPassword(r.PathValue("id"), payload)
	if err != nil {
		if errors.Is(err, service.ErrUpdatingPassword) {
			NewConflictError(w, err.Error())
			return
		}
		NewBadRequestError(w, err.Error())
		return
	}

	WriteSuccessResponse(w, http.StatusCreated)
}

func (h userHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateUser

	err := ParseBodyRequest(r, &input)
	if err != nil {
		NewBadRequestError(w, "http error parsing body request")
		return
	}

	err = h.userService.Update(r.PathValue("id"), input)
	if err != nil {
		NewBadRequestError(w, "cannot update user")
		return
	}

	WriteSuccessResponse(w, http.StatusCreated)
}

func (h userHandler) getSigned(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	payload, err := h.jwtAuth.VerifyToken(token)
	if err != nil {
		NewUnauthorizedError(w, "unauthorized user")
		return
	}

	u, err := h.userService.GetUser(payload.PublicID)
	if err != nil {
		NewBadRequestError(w, "user not found")
		return
	}

	user := response.User{
		PublicID:  u.PublicID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Document:  u.Document,
		Enabled:   u.Enabled,
		CreatedAt: u.CreatedAt,
	}

	WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
