package restutil

import (
	"context"

	"github.com/bernardinorafael/gozinho/internal/infra/rest/constant"
	"github.com/gin-gonic/gin"
)

func GetContext(c *gin.Context) (ctx context.Context) {
	ctx = c.Request.Context()
	ctx = context.WithValue(ctx, constant.AuthKey, "my-context-test")

	return ctx
}
