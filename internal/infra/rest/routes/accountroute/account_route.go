package accountroute

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine, handler *Handler, a contract.AuthService) {
	user := r.Group("/user")
	{
		user.POST("", handler.Save)
		user.GET("/:id", handler.GetByID)
		user.GET("/all", handler.GetAll)
		user.PATCH("/:id", handler.Update)
		user.DELETE("/:id", handler.Delete)
		user.PUT(":id", handler.UpdatePassword)
	}

	r.POST("/login", handler.Login)
}
