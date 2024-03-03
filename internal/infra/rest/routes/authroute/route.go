package authroute

import (
	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine, handler *UserHandler) {
	a := r.Group("/auth")
	{
		a.POST("/login", handler.Login)
		a.POST("/register", handler.Register)
	}
}
