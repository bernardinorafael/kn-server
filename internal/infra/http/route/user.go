package route

import (
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/http/middleware"
	"github.com/bernardinorafael/kn-server/internal/infra/http/response"
	"github.com/bernardinorafael/kn-server/internal/infra/http/routeutils"
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
		r.With(m.WithAuth)

		r.Get("/me", h.getSigned)
		r.Put("/{id}", h.updateUser)
	})
}

func (h userHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateUser

	err := routeutils.ParseBodyRequest(r, &input)
	if err != nil {
		routeutils.NewBadRequestError(w, "http error parsing body request")
		return
	}

	err = h.userService.Update(r.PathValue("id"), input)
	if err != nil {
		routeutils.NewBadRequestError(w, "cannot update user")
		return
	}

	routeutils.WriteSuccessResponse(w, http.StatusCreated)
}

func (h userHandler) getSigned(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	payload, err := h.jwtAuth.VerifyToken(token)
	if err != nil {
		routeutils.NewUnauthorizedError(w, "unauthorized user")
		return
	}

	u, err := h.userService.GetUser(payload.PublicID)
	if err != nil {
		routeutils.NewBadRequestError(w, "user not found")
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

	routeutils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
