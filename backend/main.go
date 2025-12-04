package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/controllers"
	"backend/entity"
	"backend/seed"
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
	if err := config.DB.AutoMigrate(
		&entity.Role{},
		&entity.User{},
		&entity.AdminProfile{},
		&entity.Application{},
		&entity.ApplicationDocument{},
		&entity.ApprovalDecision{},
		&entity.ApprovalTask{},
		&entity.FamilyInfo{},
		&entity.Major{},
		&entity.NewsPost{},
		&entity.Scholarship{},
		&entity.Screening{},
		&entity.SponsorIndustry{},
		&entity.Sponsor{},
		&entity.SponsorContact{},
		&entity.StatusNews{},
		&entity.StatusScreening{},
		&entity.Statusscholarship{},
		&entity.StudentFavNews{},
		&entity.StudentProfile{},
		&entity.Typescholarship{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	} else {
		log.Println("AutoMigrate completed")
	}

	// Seed
	seed.SeedAll(config.DB)

	// API routes
	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		api.GET("/industries", controllers.GetIndustries)
	 	api.POST("/industries", controllers.CreateIndustry)

		api.GET("/sponsors", controllers.GetSponsors)
		api.GET("/sponsors/:id", controllers.GetSponsorsByID)
		api.POST("/sponsors", controllers.CreateSponsor)
		api.PATCH("/sponsors/:id", controllers.UpdateSponsor)
		api.DELETE("/sponsors/:id", controllers.DeleteSponsor)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}
