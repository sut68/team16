package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/controllers"
	"backend/entity"
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

		r.GET("/industries", controllers.GetIndustries)
		r.POST("/industries", controllers.CreateIndustry)

		api.GET("/sponsors", controllers.GetSponsors)
		api.GET("/sponsors/:id", controllers.GetSponsorsByID)
		api.POST("/sponsors", controllers.CreateSponsor)
		api.PATCH("/sponsors/:id", controllers.UpdateSponsor)
		api.DELETE("/sponsors/:id", controllers.DeleteSponsor)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}
