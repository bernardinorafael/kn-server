package restutil

import (
	"context"

	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/gin-gonic/gin"
)

func GetContext(c *gin.Context) (ctx context.Context) {
	ctx = c.Request.Context()
	ctx = context.WithValue(ctx, auth.TokenKey, "my-context-test")

	return ctx
}
