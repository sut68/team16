package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/controllers"
	entity "backend/entity"
)

func main() {
	r := gin.Default()

	// Configure CORS
	configCORS := cors.DefaultConfig()
	configCORS.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"} // Allow frontend dev server
	configCORS.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	configCORS.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(configCORS))

	// Connect to Database
	config.ConnectDB()

	// Auto-migrate the schema
	config.DB.AutoMigrate(&entity.User{})

	// API routes
	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}
