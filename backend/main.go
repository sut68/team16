package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/controllers"
	"backend/entity"
	"backend/seed"
)

func main() {
	r := gin.New()
	r.Use(gin.LoggerWithWriter(os.Stdout))
	r.Use(gin.Recovery())

	// Configure CORS
	configCORS := cors.DefaultConfig()
	configCORS.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://localhost:5174"} // Allow frontend dev server
	configCORS.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	configCORS.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(configCORS))

	// Serve static files from the "uploads" directory
	r.Static("/uploads", "./uploads")

	// Connect to Database
	config.ConnectDB()

	if err := config.DB.Migrator().DropTable(
		&entity.Role{},
		&entity.User{},
		&entity.AdminProfile{},
		&entity.Semaster{},
		&entity.Application{},
		&entity.ApplicationScholarship{},
		&entity.ApplicationDocument{},
		&entity.ApprovalDecision{},
		&entity.ApprovalTask{},
		&entity.FamilyInfo{},
		&entity.Major{},
		&entity.NewsPost{},
		&entity.ApprovalRequirement{}, // Moved to drop before Scholarship
		&entity.Requirement{},         // Added missing Requirement
		&entity.Scholarship{},
		&entity.Featurescholarship{},
		&entity.Typefeature{},
		&entity.Screening{},
		&entity.SponsorIndustry{},
		&entity.Sponsor{},
		&entity.SponsorContact{},
		&entity.StatusNews{},
		&entity.StudentFavNews{},
		&entity.StatusScreening{},
		&entity.Statusscholarship{},
		&entity.StudentFavNews{},
		&entity.StudentProfile{},
		&entity.Typescholarship{},
		&entity.InterviewRound{},
		&entity.Interviewer{},
		&entity.Slot{},
		&entity.InterviewerSlot{},
		&entity.IntervieweBooking{},
		&entity.Location{},
		&entity.InterviewMode{},
	); err != nil {
		log.Fatalf("DropTable failed: %v", err)
	}

	//Auto-migrate the schema
	if err := config.DB.AutoMigrate(
		&entity.Role{},
		&entity.User{},
		&entity.AdminProfile{},
		&entity.Semaster{},
		&entity.Application{},
		&entity.ApplicationScholarship{},
		&entity.ApplicationDocument{},
		&entity.ApprovalDecision{},
		&entity.ApprovalTask{},
		&entity.FamilyInfo{},
		&entity.Major{},
		&entity.NewsPost{},
		&entity.Scholarship{},
		&entity.Featurescholarship{},
		&entity.Typefeature{},
		&entity.ApprovalRequirement{}, // Added missing ApprovalRequirement
		&entity.Requirement{},         // Added missing Requirement
		&entity.Screening{},
		&entity.Featurescholarship{},
		&entity.Typefeature{},
		&entity.SponsorIndustry{},
		&entity.Sponsor{},
		&entity.SponsorContact{},
		&entity.StatusNews{},
		&entity.StudentFavNews{},
		&entity.StatusScreening{},
		&entity.Statusscholarship{},
		&entity.StudentFavNews{},
		&entity.StudentProfile{},
		&entity.Typescholarship{},
		&entity.InterviewRound{},
		&entity.Interviewer{},
		&entity.Slot{},
		&entity.InterviewerSlot{},
		&entity.IntervieweBooking{},
		&entity.Location{},
		&entity.InterviewMode{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	} else {
		log.Println("AutoMigrate completed")
	}

	// Seed
	if err := seed.SeedAll(config.DB); err != nil {
		log.Fatalf("Seed failed: %v", err)
	} else {
		log.Println("Seed completed")
	}

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
		api.PATCH("/sponsors/:id/contacts", controllers.UpdateSponsorContacts)
		api.DELETE("/sponsors/:id", controllers.DeleteSponsor)

		api.GET("/students/:student_profile_id/applications", controllers.GetStudentApplications)
		// api.GET("/scholarships", controllers.GetAllScholarship)
		api.POST("/scholarship/:id/apply", controllers.ApplyForScholarship)

		//scholarship
		api.GET("/scholarship", controllers.GetAllScholarship)
		api.GET("/scholarship/:id", controllers.GetScholarshipByID)
		api.POST("/scholarship", controllers.CreateScholarship)
		api.PUT("/scholarship/:id", controllers.UpdateScholarship)
		api.DELETE("/scholarship/:id", controllers.DeleteScholarship)

		//statusscholarship
		api.GET("/statusscholarship", controllers.GetAllStatusscholarship)
		api.GET("/statusscholarship/:id", controllers.GetStatusscholarshipByID)
		api.POST("/statusscholarship", controllers.CreateStatusscholarship)
		api.PUT("/statusscholarship/:id", controllers.UpdateStatusscholarship)
		api.DELETE("/statusscholarship/:id", controllers.DeleteStatusscholarship)

		//typescholarship
		api.GET("/typescholarship", controllers.GetAllTypescholarship)
		api.GET("/typescholarship/:id", controllers.GetTypescholarshipByID)
		api.POST("/typescholarship", controllers.CreateTypescholarship)
		api.PUT("/typescholarship/:id", controllers.UpdateTypescholarship)
		api.DELETE("/typescholarship/:id", controllers.DeleteTypescholarship)

		//featurescholarship
		api.GET("/featurescholarship", controllers.GetAllFeaturescholarship)
		api.GET("/featurescholarship/:id", controllers.GetFeaturescholarshipByID)
		api.POST("/featurescholarship", controllers.CreateFeaturescholarship)
		api.PUT("/featurescholarship/:id", controllers.UpdateFeaturescholarship)
		api.DELETE("/featurescholarship/:id", controllers.DeleteFeaturescholarship)

		//typefeature
		api.GET("/typefeature", controllers.GetAllTypefeature)
		api.GET("/typefeature/:id", controllers.GetTypefeatureByID)
		api.POST("/typefeature", controllers.CreateTypefeature)
		api.PUT("/typefeature/:id", controllers.UpdateTypefeature)
		api.DELETE("/typefeature/:id", controllers.DeleteTypefeature)

		// assistance
		api.GET("/assistance", controllers.GetAllAssistance)
		api.GET("/assistance/:id", controllers.GetAssistanceByID)
		api.POST("/assistance", controllers.CreateAssistance)
		api.PUT("/assistance/:id", controllers.UpdateAssistance)
		api.DELETE("/assistance/:id", controllers.DeleteAssistance)

		api.GET("/approval-tasks", controllers.GetApprovalTasks)
		api.GET("/approval-tasks/:id", controllers.GetApprovalTaskByID)
		api.PATCH("/approval-tasks/:id", controllers.UpdateApprovalTask)
		api.DELETE("/approval-tasks/:id", controllers.DeleteApprovalTask)
		api.POST("/approval-decisions", controllers.CreateApprovalDecision)

		api.GET("/application-documents", controllers.GetApplicationDocuments)
		api.POST("/application-documents", controllers.CreateApplicationDocument)
		api.PATCH("/application-documents/:id", controllers.UpdateApplicationDocument)
		api.DELETE("/application-documents/:id", controllers.DeleteApplicationDocument)

		api.GET("/approval-requirements", controllers.GetApprovalRequirements)
		api.GET("/approval-requirements/:id", controllers.GetApprovalRequirementByID)
		api.POST("/approval-requirements", controllers.CreateApprovalRequirement)
		api.PATCH("/approval-requirements/:id", controllers.UpdateApprovalRequirement)
		api.DELETE("/approval-requirements/:id", controllers.DeleteApprovalRequirement)

		api.GET("/screening", controllers.GetAllScreenings)
		api.GET("/screening/:id", controllers.GetScreeningByID)
		api.PUT("/screening/:id", controllers.UpdateScreeningStatus)

		api.GET("/newsposts", controllers.GetAllNewsPosts)
		api.GET("/newsposts/:id", controllers.GetNewsPostByID)
		api.POST("/newsposts", controllers.CreateNewsPost)
		api.PUT("/newsposts/:id", controllers.UpdateNewsPost)
		api.DELETE("/newsposts/:id", controllers.DeleteNewsPost)

		// Interview
		api.GET("/interview-rounds", controllers.GetAllInterviewRounds)
		api.GET("/interview-rounds/:id", controllers.GetInterviewRoundByID)
		api.POST("/interview-rounds", controllers.CreateInterviewRound)
		api.PUT("/interview-rounds/:id", controllers.UpdateInterviewRound)
		api.DELETE("/interview-rounds/:id", controllers.DeleteInterviewRound)
		api.GET("/interviewers", controllers.GetAllInterviewers)
		api.POST("/interviewers", controllers.CreateInterviewer)
		api.POST("/interview-bookings", controllers.CreateInterviewBooking)
		api.GET("/students/:student_profile_id/interview-bookings", controllers.GetStudentBookings)
		api.DELETE("/interview-bookings/:id", controllers.DeleteInterviewBooking)
		api.GET("/locations", controllers.GetAllLocations)
		api.GET("/interview-modes", controllers.GetAllInterviewModes)

		// Profile (Me) - ใช้ได้ทั้ง Admin/Student Controller จะเช็คเอง
		api.GET("/profile/me", controllers.GetMyProfile)
		api.PUT("/profile/me", controllers.UpdateMyProfile)

		// User Management (Admin)
		api.GET("/roles", controllers.ListRoles)
		api.GET("/majors", controllers.ListMajors)
		api.GET("/users", controllers.ListUsers)
		api.POST("/users", controllers.CreateUser)
		api.DELETE("/users/:id", controllers.DeleteUser)
		api.PUT("/users/:id", controllers.UpdateUser) // เพิ่มเส้นทางสำหรับอัปเดตผู้ใช้

		// Student Favorite News
		api.GET("/student_favs/my_favs/:id", controllers.GetStudentFavsByStudentID)
		api.POST("/student_favs/toggle", controllers.ToggleStudentFav)

	}

	r.Run() // listen and serve on 0.0.0.0:8080
}
