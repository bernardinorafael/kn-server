package middleware

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	resterror "github.com/bernardinorafael/kn-server/internal/infra/rest/error"
	"github.com/gin-gonic/gin"
)

func Authenticate(a contract.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authentication")
		if len(token) == 0 {
			resterror.NewUnauthorizedError(c, "no Authorization (header) provided")
			c.Abort()
			return
		}

		claims, err := a.ValidateToken(c, token)
		if err != nil {
			resterror.NewInternalServerError(c, err.Error())
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
