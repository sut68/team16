package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/controllers"
	"backend/entity"
	"backend/middlewares"
	"backend/seed"
	"backend/storage"
)

func main() {
	r := gin.New()
	r.Use(gin.LoggerWithWriter(os.Stdout))
	r.Use(gin.Recovery())

	// Configure CORS
	configCORS := cors.DefaultConfig()
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		// Default for development
		corsOrigins = "http://localhost:5173,http://127.0.0.1:5173,http://localhost:5174"
	}
	configCORS.AllowOrigins = strings.Split(corsOrigins, ",")
	configCORS.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	configCORS.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "x-CSRF-Token"}
	configCORS.AllowCredentials = true
	r.Use(cors.New(configCORS))

	// Serve static files from the "uploads" directory
	r.Static("/uploads", "./uploads")

	// Connect to Database
	config.ConnectDB()

	/*
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
			&entity.Assistance{},
			&entity.Chatroom{},
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
			&entity.EvaluationScore{},
			&entity.Evaluation{},
			&entity.InterviewRoundCriterion{},
			&entity.EvaluationCriterion{},
		); err != nil {
			log.Fatalf("DropTable failed: %v", err)
		}
	*/

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
		&entity.Assistance{},
		&entity.Chatroom{},
		&entity.EvaluationCriterion{},
		&entity.InterviewRoundCriterion{},
		&entity.Evaluation{},
		&entity.EvaluationScore{},
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

	// Initialize MinIO Storage
	if err := storage.InitMinIO(); err != nil {
		log.Printf("⚠️  MinIO initialization failed: %v (uploads will use local storage)", err)
	}

	// API routes
	api := r.Group("/api")
	{
		// -----------------------------------------------------------------
		// 1. PUBLIC ROUTES (No Auth Required)
		// -----------------------------------------------------------------
		// Health Check (for Docker/Kubernetes)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.GET("/ws", controllers.WebSocketHandler)

		// ข้อมูลบริษัทผู้ให้ทุน (Sponsor)
		api.GET("/industries", controllers.GetIndustries)                         // ดึงรายการอุตสาหกรรม - หมวดหมู่บริษัท
		api.GET("/sponsors", controllers.GetSponsors)                             // ดึงรายการอุตสาหกรรม - หมวดหมู่บริษัท
		api.GET("/sponsors/:id", controllers.GetSponsorsByID)                     // ดึงข้อมูลบริษัทผู้ให้ทุน (ตาม ID)
		api.GET("/sponsors/:id/scholarships", controllers.GetSponsorScholarships) // ดึงทุนทั้งหมดของบริษัทนั้น

		// ข้อมูลทุนการศึกษา (Scholarship)
		api.GET("/scholarship", controllers.GetAllScholarship)      // ดึงรายการทุนทั้งหมด
		api.GET("/scholarship/:id", controllers.GetScholarshipByID) // ดึงรายละเอียดทุน (ตาม ID)

		// Master Data (ข้อมูลหลัก)
		api.GET("/statusscholarship", controllers.GetAllStatusscholarship)        // ดึงสถานะทุนทั้งหมด (เช่น เปิดรับ, ปิดรับ)
		api.GET("/statusscholarship/:id", controllers.GetStatusscholarshipByID)   // ดึงสถานะทุน (ตาม ID)
		api.GET("/typescholarship", controllers.GetAllTypescholarship)            // ดึงประเภททุนทั้งหมด (เช่น ทุนเต็ม, ทุนบางส่วน)
		api.GET("/typescholarship/:id", controllers.GetTypescholarshipByID)       // ดึงประเภททุน (ตาม ID)
		api.GET("/featurescholarship", controllers.GetAllFeaturescholarship)      // ดึงคุณสมบัติพิเศษทุนทั้งหมด
		api.GET("/featurescholarship/:id", controllers.GetFeaturescholarshipByID) // ดึงคุณสมบัติพิเศษทุน (ตาม ID)
		api.GET("/typefeature", controllers.GetAllTypefeature)                    // ดึงประเภทคุณสมบัติทั้งหมด
		api.GET("/typefeature/:id", controllers.GetTypefeatureByID)               // ดึงประเภทคุณสมบัติ (ตาม ID)
		api.GET("/locations", controllers.GetAllLocations)                        // ดึงสถานที่สัมภาษณ์ทั้งหมด
		api.GET("/interview-modes", controllers.GetAllInterviewModes)             // ดึงรูปแบบการสัมภาษณ์ (Onsite, Online)

		// ข่าวสาร (News)
		api.GET("/newsposts", controllers.GetAllNewsPosts)     // ดึงข่าวสารทั้งหมด
		api.GET("/newsposts/:id", controllers.GetNewsPostByID) // ดึงรายละเอียดข่าว (ตาม ID)

		// สถิติสาธารณะ (Public Statistics for Homepage)
		api.GET("/statistics/public", controllers.GetPublicStatistics) // ดึงสถิติสำหรับหน้าแรก

		// -----------------------------------------------------------------
		// 2. AUTHENTICATED ROUTES (Student / General)
		// -----------------------------------------------------------------
		authGroup := api.Group("/")
		authGroup.Use(middlewares.JWTAuth(), middlewares.CSRFMiddleware())
		{
			// การสมัครทุน (Student Application)
			authGroup.GET("/students/:student_profile_id/applications", controllers.GetStudentApplications)    // ดูรายการทุนที่ตัวเองสมัครไว้
			authGroup.POST("/scholarship/:id/apply", controllers.ApplyForScholarship)                          // สมัครทุน (ตาม Scholarship ID)
			authGroup.GET("/application-scholarships", controllers.GetAllApplicationScholarships)              // ดูการสมัครทุนทั้งหมด
			authGroup.DELETE("/application-scholarships/:id/cancel", controllers.CancelApplicationScholarship) // ยกเลิกการสมัครทุน

			// เอกสารการสมัคร (Application Documents)
			authGroup.GET("/application-documents", controllers.GetApplicationDocuments)          // ดูเอกสารที่อัปโหลดไว้
			authGroup.POST("/application-documents", controllers.CreateApplicationDocument)       // อัปโหลดเอกสาร (เช่น ใบรับรอง, สำเนาบัตร)
			authGroup.DELETE("/application-documents/:id", controllers.DeleteApplicationDocument) // ลบเอกสาร

			// ห้องแชท (Chatroom)
			authGroup.POST("/chatroom", controllers.CreateChatroom)           // สร้างห้องแชทใหม่ (เพื่อขอความช่วยเหลือ)
			authGroup.GET("/chatroom", controllers.GetAllChatroom)            // ดูห้องแชททั้งหมด
			authGroup.GET("/chatroom/:id", controllers.GetChatroomByID)       // ดูห้องแชท (ตาม ID)
			authGroup.DELETE("/chatroom/:id", controllers.DeleteChatroom)     // ลบห้องแชท
			authGroup.GET("/chatroom/my-open", controllers.GetMyOpenChatroom) // ดูห้องแชทที่ตัวเองเปิดไว้
			authGroup.GET("/chatroom/open", controllers.GetAllOpenChatrooms)  //	ดูห้องแชทที่เปิดอยู่ทั้งหมด
			authGroup.PUT("/chatroom/:id", controllers.UpdateChatroom)        // แก้ไขห้องแชท (เช่น เปลี่ยนสถานะ)

			// โปรไฟล์ (Profile)
			authGroup.GET("/profile/me", controllers.GetMyProfile)    // ดูข้อมูลโปรไฟล์ตัวเอง
			authGroup.PUT("/profile/me", controllers.UpdateMyProfile) // แก้ไขโปรไฟล์ตัวเอง

			// ข่าวที่ชอบ (Favorite News)
			authGroup.GET("/student_favs/my_favs/:id", controllers.GetStudentFavsByStudentID) // ดูข่าวที่ถูกใจ
			authGroup.POST("/student_favs/toggle", controllers.ToggleStudentFav)              // กดถูกใจ/ยกเลิกถูกใจข่าว

			// การจองสัมภาษณ์ (Interview Bookings)
			authGroup.POST("/interview-bookings", controllers.CreateInterviewBooking)                         // จองเวลาสัมภาษณ์
			authGroup.GET("/students/:student_profile_id/interview-bookings", controllers.GetStudentBookings) // ดูการจองของตัวเอง
			authGroup.DELETE("/interview-bookings/:id", controllers.DeleteInterviewBooking)                   // ยกเลิกการจอง

			// ข้อความช่วยเหลือ (Assistance/Chat Messages)
			authGroup.POST("/assistance", controllers.CreateAssistance)       // ส่งข้อความ
			authGroup.GET("/assistance", controllers.GetAllAssistance)        // ดูข้อความทั้งหมด
			authGroup.GET("/assistance/:id", controllers.GetAssistanceByID)   // ดูข้อความ (ตาม ID)
			authGroup.PUT("/assistance/:id", controllers.UpdateAssistance)    // แก้ไขข้อความ
			authGroup.DELETE("/assistance/:id", controllers.DeleteAssistance) // ลบข้อความ

			// รอบสัมภาษณ์ (Interview Rounds - Read Only)
			authGroup.GET("/interview-rounds", controllers.GetAllInterviewRounds)     // ดูรอบสัมภาษณ์ทั้งหมด (เพื่อเลือกจอง)
			authGroup.GET("/interview-rounds/:id", controllers.GetInterviewRoundByID) //	ดูรายละเอียดรอบสัมภาษณ์
		}

		// -----------------------------------------------------------------
		// 3. ADMIN ROUTES (Admin Only)
		// -----------------------------------------------------------------
		adminGroup := api.Group("/")
		adminGroup.Use(middlewares.JWTAuth(), middlewares.CSRFMiddleware(), middlewares.RequireAdmin())
		{
			// จัดการบริษัทผู้ให้ทุน (Sponsor Management)
			adminGroup.POST("/sponsors", controllers.CreateSponsor)                       // สร้างบริษัทผู้ให้ทุนใหม่
			adminGroup.PATCH("/sponsors/:id", controllers.UpdateSponsor)                  // แก้ไขข้อมูลบริษัท
			adminGroup.PATCH("/sponsors/:id/contacts", controllers.UpdateSponsorContacts) // แก้ไขข้อมูลติดต่อบริษัท
			adminGroup.DELETE("/sponsors/:id", controllers.DeleteSponsor)                 // ลบบริษัท
			adminGroup.POST("/industries", controllers.CreateIndustry)                    // สร้างหมวดหมู่อุตสาหกรรมใหม่

			// จัดการทุนการศึกษา (Scholarship Management)
			adminGroup.POST("/scholarship", controllers.CreateScholarship)       // สร้างทุนใหม่
			adminGroup.PUT("/scholarship/:id", controllers.UpdateScholarship)    // แก้ไขทุน
			adminGroup.DELETE("/scholarship/:id", controllers.DeleteScholarship) // ลบทุน

			// Master Data Management
			adminGroup.POST("/statusscholarship", controllers.CreateStatusscholarship)       // สร้างสถานะทุนใหม่
			adminGroup.PUT("/statusscholarship/:id", controllers.UpdateStatusscholarship)    // แก้ไขสถานะทุน
			adminGroup.DELETE("/statusscholarship/:id", controllers.DeleteStatusscholarship) // ลบสถานะทุน

			adminGroup.POST("/typescholarship", controllers.CreateTypescholarship)       // สร้างประเภททุนใหม่
			adminGroup.PUT("/typescholarship/:id", controllers.UpdateTypescholarship)    // แก้ไขประเภททุน
			adminGroup.DELETE("/typescholarship/:id", controllers.DeleteTypescholarship) // ลบประเภททุน

			adminGroup.POST("/featurescholarship", controllers.CreateFeaturescholarship)       // สร้างคุณสมบัติพิเศษทุน
			adminGroup.PUT("/featurescholarship/:id", controllers.UpdateFeaturescholarship)    // แก้ไขคุณสมบัติ
			adminGroup.DELETE("/featurescholarship/:id", controllers.DeleteFeaturescholarship) // ลบคุณสมบัติ

			adminGroup.POST("/typefeature", controllers.CreateTypefeature)       // สร้างประเภทคุณสมบัติ
			adminGroup.PUT("/typefeature/:id", controllers.UpdateTypefeature)    // แก้ไขประเภทคุณสมบัติ
			adminGroup.DELETE("/typefeature/:id", controllers.DeleteTypefeature) // ลบประเภทคุณสมบัติ

			// จัดการข่าวสาร (News Posts Management)
			adminGroup.POST("/newsposts", controllers.CreateNewsPost)       // สร้างข่าวใหม่
			adminGroup.PUT("/newsposts/:id", controllers.UpdateNewsPost)    // แก้ไขข่าว
			adminGroup.DELETE("/newsposts/:id", controllers.DeleteNewsPost) // ลบข่าว

			// จัดการผู้ใช้ (User Management)
			adminGroup.GET("/users", controllers.ListUsers)         // ดูรายการผู้ใช้ทั้งหมด
			adminGroup.POST("/users", controllers.CreateUser)       // สร้างผู้ใช้ใหม่
			adminGroup.DELETE("/users/:id", controllers.DeleteUser) // แก้ไขข้อมูลผู้ใช้
			adminGroup.PUT("/users/:id", controllers.UpdateUser)    // ลบผู้ใช้
			adminGroup.GET("/roles", controllers.ListRoles)         // ดูบทบาททั้งหมด (admin, student)
			adminGroup.GET("/majors", controllers.ListMajors)       //	ดูสาขาวิชาทั้งหมด

			// จัดการสัมภาษณ์ (Interview Management)
			adminGroup.GET("/scholarships/:id/qualified-applicants", controllers.GetQualifiedApplicantsForScholarship) // ดูผู้สมัครที่ผ่านเกณฑ์
			adminGroup.POST("/admin/interview-bookings", controllers.AdminCreateInterviewBooking)                      // Admin จองแทน Student
			adminGroup.POST("/admin/interview-bookings/:id/move", controllers.AdminMoveInterviewBooking)               // Admin ย้ายเวลาจอง

			adminGroup.POST("/interview-rounds", controllers.CreateInterviewRound)       // สร้างรอบสัมภาษณ์
			adminGroup.PUT("/interview-rounds/:id", controllers.UpdateInterviewRound)    // 	แก้ไขรอบสัมภาษณ์
			adminGroup.DELETE("/interview-rounds/:id", controllers.DeleteInterviewRound) // 	ลบรอบสัมภาษณ์

			adminGroup.GET("/interviewers", controllers.GetAllInterviewers) // 	ดูรายชื่อผู้สัมภาษณ์
			adminGroup.POST("/interviewers", controllers.CreateInterviewer) // 	เพิ่มผู้สัมภาษณ์

			// อนุมัติเอกสาร (Approval Tasks & Decisions)
			adminGroup.GET("/approval-tasks", controllers.GetApprovalTasks)            // 	ดูรายการเอกสารรออนุมัติ
			adminGroup.GET("/approval-tasks/:id", controllers.GetApprovalTaskByID)     // 	ดูรายละเอียดเอกสาร
			adminGroup.PATCH("/approval-tasks/:id", controllers.UpdateApprovalTask)    // 	อัปเดตสถานะเอกสาร
			adminGroup.DELETE("/approval-tasks/:id", controllers.DeleteApprovalTask)   // 	ลบเอกสาร
			adminGroup.POST("/approval-decisions", controllers.CreateApprovalDecision) // บันทึกการตัดสินใจ (อนุมัติ/ปฏิเสธ/ขอแก้ไข)

			adminGroup.GET("/approval-requirements", controllers.GetApprovalRequirements)          // ดูเงื่อนไขการอนุมัติ
			adminGroup.GET("/approval-requirements/:id", controllers.GetApprovalRequirementByID)   // ดูเงื่อนไข (ตาม ID)
			adminGroup.POST("/approval-requirements", controllers.CreateApprovalRequirement)       // สร้างเงื่อนไขใหม่
			adminGroup.DELETE("/approval-requirements/:id", controllers.DeleteApprovalRequirement) // ลบเงื่อนไข

			// คัดกรอง (Screening)
			adminGroup.GET("/screening", controllers.GetAllScreenings)          // ดูรายการคัดกรองทั้งหมด
			adminGroup.GET("/screening/:id", controllers.GetScreeningByID)      // ดูรายละเอียดการคัดกรอง
			adminGroup.PUT("/screening/:id", controllers.UpdateScreeningStatus) //	อัปเดตสถานะคัดกรอง (ผ่าน/ไม่ผ่าน)

			// ระบบประเมิน (Evaluation System)
			adminGroup.GET("/evaluation-criteria", controllers.GetAllEvaluationCriteria)         // ดูเกณฑ์การประเมินทั้งหมด
			adminGroup.GET("/evaluation-criteria/:id", controllers.GetEvaluationCriterionByID)   // 	ดูเกณฑ์ (ตาม ID)
			adminGroup.POST("/evaluation-criteria", controllers.CreateEvaluationCriterion)       // 	สร้างเกณฑ์ใหม่
			adminGroup.PATCH("/evaluation-criteria/:id", controllers.UpdateEvaluationCriterion)  // แก้ไขเกณฑ์
			adminGroup.DELETE("/evaluation-criteria/:id", controllers.DeleteEvaluationCriterion) // 	ลบเกณฑ์

			adminGroup.GET("/interview-rounds/:id/criteria", controllers.GetInterviewRoundCriteria)           // 	ดูเกณฑ์ของรอบสัมภาษณ์
			adminGroup.POST("/interview-rounds/:id/criteria", controllers.AddCriterionToInterviewRound)       // เพิ่มเกณฑ์เข้ารอบ
			adminGroup.PATCH("/interview-round-criteria/:id", controllers.UpdateInterviewRoundCriterion)      // แก้ไขเกณฑ์ในรอบ
			adminGroup.DELETE("/interview-round-criteria/:id", controllers.RemoveCriterionFromInterviewRound) // 	ลบเกณฑ์ออกจากรอบ

			adminGroup.GET("/evaluations", controllers.GetAllEvaluations)                // 	ดูการประเมินทั้งหมด
			adminGroup.GET("/evaluations/:id", controllers.GetEvaluationByID)            // 	ดูการประเมิน (ตาม ID)
			adminGroup.POST("/evaluations", controllers.CreateEvaluation)                // สร้างการประเมินใหม่
			adminGroup.PATCH("/evaluations/:id", controllers.UpdateEvaluation)           // 	แก้ไขการประเมิน
			adminGroup.DELETE("/evaluations/:id", controllers.DeleteEvaluation)          // ลบการประเมิน
			adminGroup.POST("/evaluations/:id/complete", controllers.CompleteEvaluation) // บันทึกผลการประเมินเสร็จสิ้น

			adminGroup.POST("/evaluations/:id/scores", controllers.AddEvaluationScore)     // 	เพิ่มคะแนนประเมิน
			adminGroup.PATCH("/evaluation-scores/:id", controllers.UpdateEvaluationScore)  // แก้ไขคะแนน
			adminGroup.DELETE("/evaluation-scores/:id", controllers.DeleteEvaluationScore) // ลบคะแนน
		}
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}
