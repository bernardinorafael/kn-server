package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/http/httperror"
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

func (h *authHandler) RegisterRoute(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", h.login)
	mux.HandleFunc("POST /auth/register", h.register)
	mux.HandleFunc("PATCH /auth/{id}/password", h.recoverPassword)
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	var payload dto.Login

	if r.Body == nil {
		httperror.NewBadRequestError(w, "cannot parse body")
		return
	}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		httperror.NewInternalServerError(w, "an unknown error occurred")
		return
	}

	user, err := h.authService.Login(payload.Email, payload.Password)
	if err != nil {
		httperror.NewBadRequestError(w, err.Error())
		return
	}

	token, err := h.jwtService.CreateToken(user.PublicID)
	if err != nil {
		httperror.NewInternalServerError(w, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"public_id": user.PublicID,
		"token":     token,
	})
}

func (h *authHandler) register(w http.ResponseWriter, r *http.Request) {
	var payload dto.Register

	if r.Body == nil {
		httperror.NewBadRequestError(w, "cannot parse body")
		return
	}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		httperror.NewInternalServerError(w, "an unknown error occurred")
		return
	}

	user, err := h.authService.Register(payload.Name, payload.Email, payload.Password)
	if err != nil {
		httperror.NewBadRequestError(w, err.Error())
		return
	}

	token, err := h.jwtService.CreateToken(user.PublicID)
	if err != nil {
		httperror.NewInternalServerError(w, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"public_id": user.PublicID,
		"token":     token,
	})
}

func (h *authHandler) recoverPassword(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdatePassword

	if r.Body == nil {
		httperror.NewBadRequestError(w, "cannot parse body")
		return
	}
	defer r.Body.Close()

	id := r.PathValue("id")
	parsedID, _ := strconv.ParseUint(id, 10, 64)

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		httperror.NewInternalServerError(w, "an unknown error occurred")
		return
	}

	err = h.authService.RecoverPassword(uint(parsedID), payload)
	if err != nil {
		httperror.NewInternalServerError(w, "an unknown error occurred")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
