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
	authService contract.AuthService
	l           *slog.Logger
}

func NewHandler(l *slog.Logger, authService contract.AuthService) *Handler {
	return &Handler{
		authService: authService,
		l:           l,
	}
}

func (h *Handler) RegisterRoute(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", h.login)
	mux.HandleFunc("POST /auth/register", h.register)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	// [] - parse payload
	// [] - verify user exists by email
	// [] - check password match
	// [] - generate token
	// [] - response the token
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var payload dto.Register

	// TODO: transform body verification into a helper fn and make ir better
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

	err = h.authService.Register(payload.Name, payload.Email, payload.Password)
	if err != nil {
		resterror.NewBadRequestError(w, err.Error())
		return
	}

	// TODO: transform success response into a fn
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
