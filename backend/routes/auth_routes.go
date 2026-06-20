package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	// public auth endpoints
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
	}

	// secured auth endpoints
	securedAuthGroup := router.Group("/api/auth")
	securedAuthGroup.Use(middleware.AuthMiddleware())
	{
		securedAuthGroup.GET("/me", controllers.Me)
		securedAuthGroup.POST("/logout", controllers.Logout)
	}
}
