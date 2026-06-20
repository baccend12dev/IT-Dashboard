package main

import (
	"backend/config"
	"backend/models"
	"backend/routes"
	"backend/utils"
	"os"

	_ "backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Application Knowledge Management System API
// @version         1.0
// @description     Central portal API for managing systems, servers, notes, feature requests, and documentation.
// @host            localhost:8080
// @BasePath        /
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	// Load environment variables from .env file for local development
	_ = godotenv.Load(".env", "../.env")

	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.System{})
	config.DB.AutoMigrate(&models.Server{})
	config.DB.AutoMigrate(&models.Note{})
	config.DB.AutoMigrate(&models.FeatureRequest{})
	config.DB.AutoMigrate(&models.Documentation{})
	config.DB.AutoMigrate(&models.User{})

	// Seed default Administrator user if DB is empty
	var userCount int64
	config.DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		hashedPassword, err := utils.HashPassword("admin123")
		if err == nil {
			adminUser := models.User{
				Username: "admin",
				Password: hashedPassword,
				Role:     "Administrator",
			}
			config.DB.Create(&adminUser)
		}
	}

	r := gin.Default()

	// Enable CORS
	r.Use(CORSMiddleware())

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "database connection successful",
		})
	})

	routes.SetupAuthRoutes(r)
	routes.SetupSystemRoutes(r)
	routes.SetupServerRoutes(r)
	routes.SetupNoteRoutes(r)
	routes.SetupFeatureRequestRoutes(r)
	routes.SetupDocumentationRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
