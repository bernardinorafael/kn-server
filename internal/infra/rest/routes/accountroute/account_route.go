package accountroute

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine, handler *Handler, a contract.AuthService) {
	r.POST("/user", handler.Save)
	r.GET("/user/:id", handler.GetByID)
	r.GET("/users", handler.GetAll)
	r.PATCH("/user/:id", handler.Update)
	r.DELETE("/user/:id", handler.Delete)
	r.PUT("/user/:id", handler.UpdatePassword)
	r.POST("/login", handler.Login)
}
