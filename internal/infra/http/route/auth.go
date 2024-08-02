package route

import (
	"errors"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/config"
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/http/httperr"
	"github.com/bernardinorafael/kn-server/internal/infra/http/restutil"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type authHandler struct {
	log         logger.Logger
	env         *config.Env
	authService contract.AuthService
	jwtAuth     auth.TokenAuthInterface
}

func NewAuthHandler(log logger.Logger, authService contract.AuthService, jwtAuth auth.TokenAuthInterface, env *config.Env) *authHandler {
	return &authHandler{log, env, authService, jwtAuth}
}

func (h *authHandler) RegisterRoute(r *chi.Mux) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.login)
		r.Post("/register", h.register)
		r.Patch("/{id}/password", h.recoverPassword)
	})
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	var input dto.Login

	err := restutil.ParseBody(r, &input)
	if err != nil {
		httperr.BadRequestError(w, "error parsing body request")
		return
	}

	user, err := h.authService.Login(input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredential) {
			httperr.ConflictError(w, err.Error())
			return
		}
		httperr.BadRequestError(w, err.Error())
		return
	}

	token, payload, err := h.jwtAuth.CreateAccessToken(user.PublicID, h.env.AccessTokenDuration)
	if err != nil {
		httperr.InternalServerError(w, err.Error())
		return
	}
	restutil.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"token":   token,
		"payload": payload,
	})
}

func (h *authHandler) register(w http.ResponseWriter, r *http.Request) {
	var input dto.Register

	err := restutil.ParseBody(r, &input)
	if err != nil {
		httperr.BadRequestError(w, "error parsing body request")
		return
	}

	_, err = h.authService.Register(input)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyTaken) {
			httperr.ConflictError(w, err.Error())
			return
		}
		if errors.Is(err, service.ErrDocumentAlreadyTaken) {
			httperr.ConflictError(w, err.Error())
			return
		}
		httperr.BadRequestError(w, err.Error())
		return
	}
	restutil.WriteSuccess(w, http.StatusCreated)
}

func (h *authHandler) recoverPassword(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdatePassword

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		httperr.BadRequestError(w, "error parsing body request")
		return
	}

	err = h.authService.RecoverPassword(r.PathValue("id"), payload)
	if err != nil {
		if errors.Is(err, service.ErrUpdatingPassword) {
			httperr.ConflictError(w, err.Error())
			return
		}
		httperr.BadRequestError(w, err.Error())
		return
	}
	restutil.WriteSuccess(w, http.StatusCreated)
}
