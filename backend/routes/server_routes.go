package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupServerRoutes(router *gin.Engine) {
	serverGroup := router.Group("/api/servers")
	serverGroup.Use(middleware.AuthMiddleware())
	{
		serverGroup.GET("/", controllers.GetServers)
		serverGroup.GET("/:id", controllers.GetServerByID)

		writeGroup := serverGroup.Group("")
		writeGroup.Use(middleware.RoleMiddleware("Administrator", "Developer"))
		{
			writeGroup.POST("/", controllers.CreateServer)
			writeGroup.PUT("/:id", controllers.UpdateServer)
			writeGroup.DELETE("/:id", controllers.DeleteServer)
		}
	}
}
