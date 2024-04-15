package userroute

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/middleware"
	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine, handler *UserHandler, js contract.JWT) {
	user := r.Group("/")
	user.Use(middleware.Authenticate(js))
	{
		user.GET("/users", handler.GetManyUsers)
		user.GET("/user/:id", handler.GetUser)
	}
}
