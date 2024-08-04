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
	log             logger.Logger
	env             *config.Env
	authService     contract.AuthService
	notifierService contract.SMSNotifier
	jwtAuth         auth.TokenAuthInterface
}

func NewAuthHandler(
	log logger.Logger,
	authService contract.AuthService,
	notifierService contract.SMSNotifier,
	jwtAuth auth.TokenAuthInterface,
	env *config.Env,
) *authHandler {
	return &authHandler{
		log:             log,
		env:             env,
		authService:     authService,
		notifierService: notifierService,
		jwtAuth:         jwtAuth,
	}
}

func (h authHandler) RegisterRoute(r *chi.Mux) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.login)
		r.Post("/register", h.register)
		r.Post("/notify", h.notify)
		r.Post("/verify", h.verify)
	})
}

func (h authHandler) verify(w http.ResponseWriter, r *http.Request) {
	var input dto.VerifySMS

	if err := ParseBodyRequest(r, &input); err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	status, err := h.notifierService.Confirm(input.Code, input.Phone)
	if err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	switch {
	case status == "pending":
		NewBadRequestError(w, "invalid code")
		return
	case status == "canceled":
		NewBadRequestError(w, "verify operation canceled")
		return
	case status == "max_attempts_reached":
		NewBadRequestError(w, "max attempts reached")
		return
	case status == "expired":
		NewBadRequestError(w, "the verification has expired")
		return
	}

	WriteSuccessResponse(w, http.StatusOK)
}

func (h authHandler) notify(w http.ResponseWriter, r *http.Request) {
	var input dto.NotifySMS

	if err := ParseBodyRequest(r, &input); err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	if err := h.notifierService.Notify(input.Phone); err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	WriteSuccessResponse(w, http.StatusOK)
}

func (h authHandler) login(w http.ResponseWriter, r *http.Request) {
	var input dto.Login

	if err := ParseBodyRequest(r, &input); err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	user, err := h.authService.Login(input)
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
	var input dto.Register

	if err := ParseBodyRequest(r, &input); err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	_, err := h.authService.Register(input)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyTaken) {
			NewConflictError(w, err.Error())
			return
		}
		if errors.Is(err, service.ErrDocumentAlreadyTaken) {
			NewConflictError(w, err.Error())
			return
		}
		NewBadRequestError(w, err.Error())
		return
	}

	WriteSuccessResponse(w, http.StatusCreated)
}
