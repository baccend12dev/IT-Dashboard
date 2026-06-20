package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupSystemRoutes(router *gin.Engine) {
	systemGroup := router.Group("/api/systems")
	systemGroup.Use(middleware.AuthMiddleware())
	{
		systemGroup.GET("/", controllers.GetSystems)
		systemGroup.GET("/:id", controllers.GetSystemByID)

		writeGroup := systemGroup.Group("")
		writeGroup.Use(middleware.RoleMiddleware("Administrator", "Developer"))
		{
			writeGroup.POST("/", controllers.CreateSystem)
			writeGroup.PUT("/:id", controllers.UpdateSystem)
			writeGroup.DELETE("/:id", controllers.DeleteSystem)
		}
	}
}
