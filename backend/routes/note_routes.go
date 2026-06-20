package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupNoteRoutes(router *gin.Engine) {
	noteRoutes := router.Group("/api/systems/:id/notes")
	noteRoutes.Use(middleware.AuthMiddleware())
	{
		noteRoutes.GET("", controllers.GetNotesBySystemID)

		writeGroup := noteRoutes.Group("")
		writeGroup.Use(middleware.RoleMiddleware("Administrator", "Developer"))
		{
			writeGroup.POST("", controllers.CreateNote)
		}
	}

	noteIndividual := router.Group("/api/notes/:id")
	noteIndividual.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("Administrator", "Developer"))
	{
		noteIndividual.PUT("", controllers.UpdateNote)
		noteIndividual.DELETE("", controllers.DeleteNote)
	}
}
