package middleware

import (
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	resterror "github.com/bernardinorafael/kn-server/internal/infra/rest/error"
	"github.com/gin-gonic/gin"
)

func WithAuthentication(a auth.Authentication) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(string(auth.TokenKey))
		if len(token) == 0 {
			resterror.NewUnauthorizedRequestError(c, "access token is required")
			return
		}

		payload, err := a.ValidateToken(c, token)
		if err != nil {
			resterror.NewUnauthorizedRequestError(c, "auth failed")
			return
		}

		c.Set(string(auth.UserIDKey), payload.UserID)
		c.Next()
	}
}
