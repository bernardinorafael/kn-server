package middleware

import (
	"github.com/bernardinorafael/gozinho/internal/application/contract"
	"github.com/bernardinorafael/gozinho/internal/infra"
	resterror "github.com/bernardinorafael/gozinho/internal/infra/rest/error"
	"github.com/gin-gonic/gin"
)

func WithAuthentication(authToken contract.Authentication) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(string(infra.TokenKey))
		if len(token) == 0 {
			resterror.NewUnauthorizedRequestError(c, "authentication required")
			return
		}

		payload, err := authToken.VerifyToken(c, token)
		if err != nil {
			resterror.NewUnauthorizedRequestError(c, "authentication required")
		}

		c.Set(string(infra.UserIDKey), payload.UserID)
		c.Next()
	}
}
