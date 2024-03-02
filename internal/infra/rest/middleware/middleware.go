package middleware

import (
	"net/http"
	"strings"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	resterr "github.com/bernardinorafael/kn-server/internal/infra/rest/error"
	"github.com/gin-gonic/gin"
)

func Authenticate(a contract.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, resterr.HttpResponseErr{
				Err:     "bad_request",
				Message: "authorization header not provided",
				Code:    http.StatusBadRequest,
			})
			return
		}

		p := strings.Split(header, " ")
		if len(p) != 2 || p[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, resterr.HttpResponseErr{
				Err:     "unauthorized",
				Message: "invalid access_token/Authorization bearer header",
				Code:    http.StatusBadRequest,
			})
			return
		}

		token := p[1]
		claims, err := a.ValidateToken(c, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resterr.HttpResponseErr{
				Err:     "internal_server_error",
				Message: "failed to validate access_token",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
