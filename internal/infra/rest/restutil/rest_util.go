package restutil

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Key string

const AuthKey Key = "user_id"

func GetContext(c *gin.Context) (ctx context.Context) {
	ctx = c.Request.Context()

	if userId, ok := c.Get(string(AuthKey)); ok {
		ctx = context.WithValue(ctx, AuthKey, userId)
	}
	return ctx
}
