package route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/config"
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/http/error"
	"github.com/bernardinorafael/kn-server/internal/infra/http/restutil"
	"github.com/bernardinorafael/kn-server/internal/infra/http/server"
)

type authHandler struct {
	log         *slog.Logger
	env         *config.Env
	authService contract.AuthService
	jwtAuth     auth.TokenAuthInterface
}

func NewAuthHandler(log *slog.Logger, authService contract.AuthService, jwtAuth auth.TokenAuthInterface, env *config.Env) *authHandler {
	return &authHandler{log, env, authService, jwtAuth}
}

func (h *authHandler) RegisterRoute(s *server.Server) {
	s.Group(func(s *server.Server) {
		s.Post("/auth/login", h.login)
		s.Post("/auth/register", h.register)
		s.Patch("/auth/{id}/password", h.recoverPassword)
	})
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	var input dto.Login

	err := restutil.ParseBody(r, &input)
	if err != nil {
		error.NewBadRequestError(w, "error parsing body request")
		return
	}

	user, err := h.authService.Login(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredential) {
			error.NewConflictError(w, err.Error())
			return
		}
		error.NewBadRequestError(w, err.Error())
		return
	}

	token, payload, err := h.jwtAuth.CreateAccessToken(user.PublicID, h.env.AccessTokenDuration)
	if err != nil {
		error.NewInternalServerError(w, err.Error())
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
		error.NewBadRequestError(w, "error parsing body request")
		return
	}

	_, err = h.authService.Register(input.Name, input.Email, input.Password, input.Document)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyTaken) {
			error.NewConflictError(w, err.Error())
			return
		}
		if errors.Is(err, service.ErrDocumentAlreadyTaken) {
			error.NewConflictError(w, err.Error())
			return
		}
		error.NewBadRequestError(w, err.Error())
		return
	}
	restutil.WriteSuccess(w, http.StatusCreated)
}

func (h *authHandler) recoverPassword(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdatePassword

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		error.NewBadRequestError(w, "error parsing body request")
		return
	}

	publicID := r.PathValue("id")

	err = h.authService.RecoverPassword(publicID, payload)
	if err != nil {
		if errors.Is(err, service.ErrUpdatingPassword) {
			error.NewConflictError(w, err.Error())
			return
		}
		error.NewBadRequestError(w, err.Error())
		return
	}
	restutil.WriteSuccess(w, http.StatusCreated)
}
