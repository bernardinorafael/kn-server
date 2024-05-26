package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/http/resterror"
)

type Handler struct {
	l           *slog.Logger
	authService contract.AuthService
	jwtService  contract.JWTService
}

func NewHandler(l *slog.Logger, authService contract.AuthService, jwtService contract.JWTService) *Handler {
	return &Handler{
		l:           l,
		authService: authService,
		jwtService:  jwtService,
	}
}

func (h *Handler) RegisterRoute(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", h.login)
	mux.HandleFunc("POST /auth/register", h.register)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var payload dto.Login

	if r.Body == nil {
		resterror.NewBadRequestError(w, "cannot parse body")
		return
	}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		resterror.NewInternalServerError(w, "an unknown error occurred")
		return
	}

	user, err := h.authService.Login(payload.Email, payload.Password)
	if err != nil {
		resterror.NewBadRequestError(w, err.Error())
		return
	}

	token, err := h.jwtService.CreateToken(user.ID)
	if err != nil {
		resterror.NewInternalServerError(w, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": user.ID,
		"token":   token,
	})
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var payload dto.Register

	if r.Body == nil {
		resterror.NewBadRequestError(w, "cannot parse body")
		return
	}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		resterror.NewInternalServerError(w, "an unknown error occurred")
		return
	}

	user, err := h.authService.Register(payload.Name, payload.Email, payload.Password)
	if err != nil {
		resterror.NewBadRequestError(w, err.Error())
		return
	}

	token, err := h.jwtService.CreateToken(user.ID)
	if err != nil {
		resterror.NewInternalServerError(w, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": user.ID,
		"token":   token,
	})
}
