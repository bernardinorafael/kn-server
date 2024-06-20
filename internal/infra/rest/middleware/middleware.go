package middleware

import (
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/error"
)

type middleware struct {
	log     *slog.Logger
	jwtAuth auth.TokenAuthInterface
}

func New(jwtAuth auth.TokenAuthInterface, log *slog.Logger) *middleware {
	return &middleware{log, jwtAuth}
}

func (m *middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if len(accessToken) == 0 {
			m.log.Error("authorization header not found")
			error.NewUnauthorizedError(w, "unauthorized user")
			return
		}

		_, err := m.jwtAuth.VerifyToken(accessToken)
		if err != nil {
			m.log.Error("invalid token")
			error.NewUnauthorizedError(w, "unauthorized user")
			return
		}
		next.ServeHTTP(w, r)
	})
}
