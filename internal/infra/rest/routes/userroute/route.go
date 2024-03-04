package userroute

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/middleware"
	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine, handler *UserHandler, a contract.JWTService) {
	u := r.Group("/")
	u.Use(middleware.Authenticate(a))
	{
		u.GET("/users", handler.GetAccounts)
		u.GET("/user/:id", handler.GetUser)
		u.PATCH("/user/:id", handler.UpdateUser)
		u.DELETE("/user/:id", handler.DeleteUser)
	}
}
