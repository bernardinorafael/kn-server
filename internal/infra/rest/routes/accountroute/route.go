package accountroute

import "github.com/gin-gonic/gin"

func Start(router *gin.Engine, ctrl *Handler) {
	router.POST("/user", ctrl.Save)
	router.GET("/user/:id", ctrl.GetByID)
	router.GET("/users", ctrl.GetAll)
	router.PATCH("/user/:id", ctrl.Update)
	router.DELETE("/user/:id", ctrl.Delete)
	router.PUT("/user/:id", ctrl.UpdatePassword)
}
