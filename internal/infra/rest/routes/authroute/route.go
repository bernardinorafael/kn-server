package authroute

import (
	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine, handler *UserHandler) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
	}
}
