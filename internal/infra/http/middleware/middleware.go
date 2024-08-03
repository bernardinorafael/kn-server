package middleware

import (
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/http/routeutils"
	"github.com/bernardinorafael/kn-server/pkg/logger"
)

type middleware struct {
	log     logger.Logger
	jwtAuth auth.TokenAuthInterface
}

func NewWithAuth(jwtAuth auth.TokenAuthInterface, log logger.Logger) *middleware {
	return &middleware{log, jwtAuth}
}

func (m *middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if len(accessToken) == 0 {
			m.log.Error("unauthorized", "error", "authorization header not found")
			routeutils.NewUnauthorizedError(w, "access token not found")
			return
		}

		_, err := m.jwtAuth.VerifyToken(accessToken)
		if err != nil {
			m.log.Error("unauthorized access attempt", "error", "invalid access token")
			routeutils.NewUnauthorizedError(w, "unauthorized user")
			return
		}

		next.ServeHTTP(w, r)
	})
}
