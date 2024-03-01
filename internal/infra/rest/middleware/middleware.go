package middleware

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	resterror "github.com/bernardinorafael/kn-server/internal/infra/rest/error"
	"github.com/gin-gonic/gin"
)

func Authenticate(a contract.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(string(auth.TokenKey))
		if len(token) == 0 {
			resterror.NewUnauthorizedRequestError(c, "access token not found")
			c.Abort()
			return
		}

		claims, err := a.ValidateToken(c, token)
		if err != nil {
			resterror.NewUnauthorizedRequestError(c, "unauthorized")
			c.Abort()
			return
		}

		c.Set(string(auth.UserIDKey), claims.UserID)
		c.Next()
	}
}
