package route

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
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

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		restutil.NewBadRequestError(w, "cannot parse body")
		return
	}

	if payload.Email == "" || payload.Password == "" {
		restutil.NewBadRequestError(w, "missing body request")
		return
	}

	user, err := h.authService.Login(payload.Email, payload.Password)
	if err != nil {
		restutil.NewBadRequestError(w, err.Error())
		return
	}

	token, err := h.jwtService.CreateToken(user.PublicID)
	if err != nil {
		restutil.NewInternalServerError(w, err.Error())
		return
	}

	restutil.WriteSuccess(w, http.StatusCreated, map[string]interface{}{
		"public_id": user.PublicID,
		"token":     token,
	})
}

func (h *authHandler) register(w http.ResponseWriter, r *http.Request) {
	var payload dto.Register

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		restutil.NewBadRequestError(w, "cannot parse body")
		return
	}

	user, err := h.authService.Register(payload.Name, payload.Email, payload.Password)
	if err != nil {
		restutil.NewBadRequestError(w, err.Error())
		return
	}

	token, err := h.jwtService.CreateToken(user.PublicID)
	if err != nil {
		restutil.NewInternalServerError(w, err.Error())
		return
	}

	restutil.WriteSuccess(w, http.StatusCreated, map[string]interface{}{
		"public_id": user.PublicID,
		"token":     token,
	})
}

func (h *authHandler) recoverPassword(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdatePassword

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		restutil.NewBadRequestError(w, "cannot parse body")
		return
	}

	id := r.PathValue("id")
	parsedID, _ := strconv.ParseInt(id, 10, 8)

	err = h.authService.RecoverPassword(int(parsedID), payload)
	if err != nil {
		restutil.NewInternalServerError(w, "an unknown error occurred")
		return
	}

	restutil.WriteSuccess(w, http.StatusCreated, map[string]interface{}{
		"message": "success",
	})
}
