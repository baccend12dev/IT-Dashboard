package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupDocumentationRoutes(router *gin.Engine) {
	// routes for documentation under a system
	documentationRoutes := router.Group("/api/systems/:id/documentations")
	documentationRoutes.Use(middleware.AuthMiddleware())
	{
		documentationRoutes.GET("", controllers.GetDocumentationsBySystemID)

		writeGroup := documentationRoutes.Group("")
		writeGroup.Use(middleware.RoleMiddleware("Administrator", "Developer"))
		{
			writeGroup.POST("", controllers.CreateDocumentation)
		}
	}

	// routes for single documentation operations
	router.GET("/api/documentations/:id", middleware.AuthMiddleware(), controllers.GetDocumentationByID)

	docIndividual := router.Group("/api/documentations/:id")
	docIndividual.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("Administrator", "Developer"))
	{
		docIndividual.PUT("", controllers.UpdateDocumentation)
		docIndividual.DELETE("", controllers.DeleteDocumentation)
	}
}
