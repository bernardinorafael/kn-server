package middleware

import (
	"strings"

	httperr "github.com/bernardinorafael/kn-server/helper/error"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/gin-gonic/gin"
)

func Authenticate(a contract.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if header == "" {
			httperr.NewNotFoundError(c, "missing Authorization header")
			c.Abort()
			return
		}

		p := strings.Split(header, " ")
		if len(p) != 2 || p[0] != "" {
			httperr.NewBadRequestError(c, "invalid Authorization header")
			c.Abort()
			return
		}

		token := p[1]
		claims, err := a.ValidateToken(token)
		if err != nil {
			httperr.NewInternalServerError(c, "failed to validate access token")
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
