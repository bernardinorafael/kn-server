package accountroute

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/middleware"
	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine, handler *Handler, a contract.AuthService) {
	u := r.Group("/")
	u.Use(middleware.Authenticate(a))
	{
		u.GET("/user/:id", handler.GetByID)
		u.GET("/users", handler.GetAll)
		u.PATCH("/user/:id", handler.Update)
		u.DELETE("/user/:id", handler.Delete)
		u.PUT("/user/:id", handler.UpdatePassword)
	}
	r.POST("/login", handler.Login)
	r.POST("/register", handler.Save)
}
