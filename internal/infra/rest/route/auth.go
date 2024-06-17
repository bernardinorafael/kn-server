package route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/server"
)

type authHandler struct {
	log         *slog.Logger
	authService contract.AuthService
	jwtService  contract.JWTService
}

func NewAuthHandler(log *slog.Logger, authService contract.AuthService, jwtService contract.JWTService) *authHandler {
	return &authHandler{
		log:         log,
		authService: authService,
		jwtService:  jwtService,
	}
}

func (h *authHandler) RegisterRoute(s *server.Server) {
	s.Group(func(s *server.Server) {
		s.Post("/auth/login", h.login)
		s.Post("/auth/register", h.register)
		s.Patch("/auth/{id}/password", h.recoverPassword)
	})
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	var payload dto.Login

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		restutil.NewBadRequestError(w, err.Error())
		return
	}

	user, err := h.authService.Login(payload.Email, payload.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredential) {
			restutil.NewConflictError(w, err.Error())
			return
		}
		restutil.NewBadRequestError(w, err.Error())
		return
	}

	token, err := h.jwtService.CreateToken(user.PublicID)
	if err != nil {
		restutil.NewInternalServerError(w, err.Error())
		return
	}
	restutil.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"public_id": user.PublicID,
		"token":     token,
	})
}

func (h *authHandler) register(w http.ResponseWriter, r *http.Request) {
	var payload dto.Register

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		restutil.NewBadRequestError(w, err.Error())
		return
	}

	_, err = h.authService.Register(payload.Name, payload.Email, payload.Password, payload.Document)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyTaken) {
			restutil.NewConflictError(w, err.Error())
			return
		}
		if errors.Is(err, service.ErrDocumentAlreadyTaken) {
			restutil.NewConflictError(w, err.Error())
			return
		}
		restutil.NewBadRequestError(w, err.Error())
		return
	}
	restutil.WriteSuccess(w, http.StatusCreated)
}

func (h *authHandler) recoverPassword(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdatePassword

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		restutil.NewBadRequestError(w, err.Error())
		return
	}

	publicID := r.PathValue("id")

	err = h.authService.RecoverPassword(publicID, payload)
	if err != nil {
		if errors.Is(err, service.ErrUpdatingPassword) {
			restutil.NewConflictError(w, err.Error())
			return
		}
		restutil.NewBadRequestError(w, err.Error())
		return
	}
	restutil.WriteSuccess(w, http.StatusCreated)
}
