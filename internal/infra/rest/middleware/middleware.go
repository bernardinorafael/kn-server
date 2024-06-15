package middleware

import (
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
)

type middleware struct {
	jwtService contract.JWTService
	log        *slog.Logger
}

func New(jwtService contract.JWTService, log *slog.Logger) *middleware {
	return &middleware{
		jwtService: jwtService,
		log:        log,
	}
}

// TODO: implement more robust token validation
func (m *middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if len(accessToken) == 0 {
			m.log.Error("authorization header not found")
			restutil.NewUnauthorizedError(w, "unauthorized user")
			return
		}

		_, err := m.jwtService.ValidateToken(accessToken)
		if err != nil {
			m.log.Error("invalid token")
			restutil.NewUnauthorizedError(w, "unauthorized user")
			return
		}

		next.ServeHTTP(w, r)
	})
}
