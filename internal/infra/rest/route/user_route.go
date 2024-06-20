package route

import (
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/error"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/response"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/server"
)

type userHandler struct {
	userService contract.UserService
	jwtAuth     auth.TokenAuthInterface
}

func NewUserHandler(userService contract.UserService, jwtAuth auth.TokenAuthInterface) *userHandler {
	return &userHandler{userService, jwtAuth}
}

func (h *userHandler) RegisterRoute(s *server.Server) {
	s.Get("/me", h.getSigned)
}

func (h *userHandler) getSigned(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	payload, err := h.jwtAuth.VerifyToken(token)
	if err != nil {
		error.NewUnauthorizedError(w, "unauthorized user")
		return
	}

	u, err := h.userService.GetUser(payload.PublicID)
	if err != nil {
		error.NewBadRequestError(w, "user not found")
		return
	}

	user := response.User{
		PublicID:  u.PublicID,
		Name:      u.Name,
		Email:     u.Email,
		Document:  u.Document,
		Enabled:   u.Enabled,
		CreatedAt: u.CreatedAt,
	}
	restutil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
