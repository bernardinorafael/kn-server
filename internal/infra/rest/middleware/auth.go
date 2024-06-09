package middleware

import (
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
)

type authMiddleware struct {
	jwtService contract.JWTService
}

func NewWithAuth(jwtService contract.JWTService) *authMiddleware {
	return &authMiddleware{jwtService: jwtService}
}

func (m *authMiddleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if len(accessToken) == 0 {
			restutil.NewUnauthorizedError(w, "unauthorized user")
		}

		_, err := m.jwtService.ValidateToken(accessToken)
		if err != nil {
			restutil.NewUnauthorizedError(w, "unauthorized user")
			return
		}

		next.ServeHTTP(w, r)
	})

}
