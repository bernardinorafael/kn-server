package route

import (
	"errors"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/config"
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	. "github.com/bernardinorafael/kn-server/internal/infra/http/routeutils"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type authHandler struct {
	log         logger.Logger
	env         *config.Env
	authService contract.AuthService
	jwtAuth     auth.TokenAuthInterface
}

func NewAuthHandler(
	log logger.Logger,
	authService contract.AuthService,
	jwtAuth auth.TokenAuthInterface,
	env *config.Env,
) *authHandler {
	return &authHandler{
		log:         log,
		env:         env,
		authService: authService,
		jwtAuth:     jwtAuth,
	}
}

func (h authHandler) RegisterRoute(r *chi.Mux) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.login)
		r.Post("/register", h.register)
		r.Post("/{id}/notify-validation", h.notifyValidation)
		r.Post("/login-otp", h.notifyLoginOTP)
		r.Post("/verify-otp", h.verifyLoginOTP)
	})
}

func (h authHandler) notifyValidation(w http.ResponseWriter, r *http.Request) {
	err := h.authService.NotifyValidationByEmail(r.PathValue("id"))
	if err != nil {
		NewBadRequestError(w, "error notifying user")
		return
	}
	WriteSuccessResponse(w, http.StatusOK)
}

func (h authHandler) verifyLoginOTP(w http.ResponseWriter, r *http.Request) {
	var body dto.LoginOTP

	if err := ReadRequestBody(w, r, &body); err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	user, err := h.authService.LoginOTP(body)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			NewNotFoundError(w, err.Error())
			return
		}
		NewBadRequestError(w, err.Error())
		return
	}

	token, payload, err := h.jwtAuth.CreateAccessToken(user.PublicID, h.env.AccessTokenDuration)
	if err != nil {
		NewInternalServerError(w, err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"token":   token,
		"payload": payload,
	})
}

func (h authHandler) notifyLoginOTP(w http.ResponseWriter, r *http.Request) {
	var body dto.NotifySMS

	if err := ReadRequestBody(w, r, &body); err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	if err := h.authService.NotifyLoginOTP(body); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			NewNotFoundError(w, err.Error())
			return
		}
		NewBadRequestError(w, err.Error())
		return
	}

	WriteSuccessResponse(w, http.StatusOK)
}

func (h authHandler) login(w http.ResponseWriter, r *http.Request) {
	var body dto.Login

	if err := ReadRequestBody(w, r, &body); err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	user, err := h.authService.Login(body)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredential) {
			NewConflictError(w, err.Error())
			return
		}
		NewBadRequestError(w, err.Error())
		return
	}

	token, payload, err := h.jwtAuth.CreateAccessToken(user.PublicID, h.env.AccessTokenDuration)
	if err != nil {
		NewInternalServerError(w, err.Error())
		return
	}

	WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"token":   token,
		"payload": payload,
	})
}

func (h authHandler) register(w http.ResponseWriter, r *http.Request) {
	var body dto.Register

	if err := ReadRequestBody(w, r, &body); err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	err := h.authService.Register(body)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyTaken) {
			NewConflictError(w, err.Error())
			return
		}
		NewBadRequestError(w, err.Error())
		return
	}

	WriteSuccessResponse(w, http.StatusCreated)
}
