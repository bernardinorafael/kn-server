package middleware

import (
	"net/http"

	. "github.com/bernardinorafael/kn-server/internal/infra/http/routeutils"
)

func WithRecoverPanic(done http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				NewInternalServerError(w, "internal server error")
			}
		}()
		done.ServeHTTP(w, r)
	})
}
