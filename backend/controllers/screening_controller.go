package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/entity"
	"backend/services"
	"backend/ws"
)

func GetAllScreenings(c *gin.Context) {
	var screenings []entity.Screening

	// Only preload fields required for the list view
	err := config.DB.
		Preload("StatusScreening").
		Preload("ApplicationScholarship.Application.StudentProfile").
		Preload("ApplicationScholarship.Scholarship").
		Preload("ApplicationScholarship.Scholarship.Semaster").
		Order("created_at desc").
		Find(&screenings).Error

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Create a lightweight response structure
	type StudentProfileDTO struct {
		FirstNameTH string `json:"first_name_th"`
		LastNameTH  string `json:"last_name_th"`
	}
	type ApplicationDTO struct {
		StudentProfile StudentProfileDTO `json:"student_profile"`
	}
	type SemasterDTO struct {
		Term         string `json:"Term"` // Match frontend expectation
		AcademicYear string `json:"AcademicYear"`
		Round        string `json:"round"`
	}
	type ScholarshipDTO struct {
		ScholarshipName string      `json:"scholarship_name"`
		Semaster        SemasterDTO `json:"semaster"`
	}
	type AppScholarshipDTO struct {
		Application ApplicationDTO `json:"application"`
		Scholarship ScholarshipDTO `json:"scholarship"`
	}
	type ScreeningDTO struct {
		ID                     uint              `json:"ID"`
		CreatedAt              string            `json:"CreatedAt"`
		StatusScreeningID      uint              `json:"status_screening_id"`
		RejectionReason        *string           `json:"rejection_reason"`
		ApplicationScholarship AppScholarshipDTO `json:"application_scholarship"`
	}

	var response []ScreeningDTO
	for _, s := range screenings {
		appSch := s.ApplicationScholarship
		sch := appSch.Scholarship
		sem := sch.Semaster
		app := appSch.Application
		student := app.StudentProfile

		dto := ScreeningDTO{
			ID:                s.ID,
			CreatedAt:         s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			StatusScreeningID: s.StatusScreeningID,
			RejectionReason:   s.RejectionReason,
			ApplicationScholarship: AppScholarshipDTO{
				Scholarship: ScholarshipDTO{
					ScholarshipName: sch.ScholarshipName,
					Semaster: SemasterDTO{
						Term:         sem.Term,
						AcademicYear: sem.AcademicYear,
						Round:        sem.Round,
					},
				},
				Application: ApplicationDTO{
					StudentProfile: StudentProfileDTO{
						FirstNameTH: student.FirstNameTH,
						LastNameTH:  student.LastNameTH,
					},
				},
			},
		}
		response = append(response, dto)
	}

	c.JSON(200, gin.H{"data": response})
}

func GetScreeningByID(c *gin.Context) {
	var screening entity.Screening
	id := c.Param("id")

	// Load screening + related
	if err := config.DB.
		Preload("StatusScreening").
		Preload("AdminProfile").
		Preload("ApplicationScholarship.Application.StudentProfile").
		Preload("ApplicationScholarship.Application.StudentProfile.FamilyInfo").
		Preload("ApplicationScholarship.Scholarship").
		Preload("ApplicationScholarship.Scholarship.Semaster").
		First(&screening, "id = ?", id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Data not found"})
		return
	}

	// Load features separately
	var features []entity.Featurescholarship
	if err := config.DB.
		Preload("Typefeature").
		Where("scholarship_id = ?", screening.ApplicationScholarship.ScholarshipID).
		Find(&features).Error; err != nil {
		features = []entity.Featurescholarship{}
	}

	// Wrap JSON manually without Sponsor
	response := gin.H{
		"data": gin.H{
			"ID":                         screening.ID,
			"CreatedAt":                  screening.CreatedAt,
			"UpdatedAt":                  screening.UpdatedAt,
			"DeletedAt":                  screening.DeletedAt,
			"admin_profile_id":           screening.AdminProfileID,
			"admin_profile":              screening.AdminProfile,
			"application_scholarship_id": screening.ApplicationScholarshipID,
			"application_scholarship": gin.H{
				"ID":             screening.ApplicationScholarship.ID,
				"status":         screening.ApplicationScholarship.Status,
				"application_id": screening.ApplicationScholarship.ApplicationID,
				"application":    screening.ApplicationScholarship.Application,
				"scholarship_id": screening.ApplicationScholarship.ScholarshipID,
				"scholarship": gin.H{
					"ID":                  screening.ApplicationScholarship.Scholarship.ID,
					"scholarship_name":    screening.ApplicationScholarship.Scholarship.ScholarshipName,
					"description":         screening.ApplicationScholarship.Scholarship.Description,
					"open_date":           screening.ApplicationScholarship.Scholarship.OpenDate,
					"close_date":          screening.ApplicationScholarship.Scholarship.CloseDate,
					"semaster":            screening.ApplicationScholarship.Scholarship.Semaster,
					"featurescholarships": features, // attach แยก
				},
			},
			"status_screening_id": screening.StatusScreeningID,
			"status_screening":    screening.StatusScreening,
			"rejection_reason":    screening.RejectionReason,
		},
	}

	c.JSON(200, response)
}

func UpdateScreeningStatus(c *gin.Context) {
	var screening entity.Screening
	id := c.Param("id")

	if err := config.DB.First(&screening, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Screening not found"})
		return
	}

	var input struct {
		StatusScreeningID uint    `json:"status_screening_id" binding:"required"`
		RejectionReason   *string `json:"rejection_reason"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Fetch current admin info from token
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var admin entity.AdminProfile
	if err := config.DB.Where("user_id = ?", userID).First(&admin).Error; err != nil {
		c.JSON(400, gin.H{"error": "Admin profile not found"})
		return
	}

	const (
		PASS = 2
		FAIL = 3
	)

	if input.StatusScreeningID == FAIL && input.RejectionReason == nil {
		c.JSON(400, gin.H{"error": "Rejection reason is required"})
		return
	}

	if input.StatusScreeningID == PASS {
		input.RejectionReason = nil
	}

	tx := config.DB.Begin()

	if err := tx.Model(&screening).Updates(map[string]interface{}{
		"status_screening_id": input.StatusScreeningID,
		"rejection_reason":    input.RejectionReason,
		"admin_profile_id":    admin.ID,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update ApplicationScholarship status based on screening result
	switch input.StatusScreeningID {
	case PASS:
		// When screening passes, set status to empty (ready for document upload - step 3)
		if err := tx.Model(&entity.ApplicationScholarship{}).
			Where("id = ?", screening.ApplicationScholarshipID).
			Update("status", "").Error; err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"error": "Failed to update application status"})
			return
		}
	case FAIL:
		// When screening fails, set status to rejected
		if err := tx.Model(&entity.ApplicationScholarship{}).
			Where("id = ?", screening.ApplicationScholarshipID).
			Update("status", "rejected").Error; err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"error": "Failed to update application status"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Transaction commit failed"})
		return
	}

	// Broadcast to WebSocket for real-time update on StudentDocument
	if ws.ApprovalHubInstance != nil {
		log.Printf("📢 Broadcasting screening status updated: ID=%d, Status=%d", screening.ID, input.StatusScreeningID)
		ws.ApprovalHubInstance.BroadcastUpdate("screening_status_updated", gin.H{
			"screening_id":               screening.ID,
			"application_scholarship_id": screening.ApplicationScholarshipID,
			"status_screening_id":        input.StatusScreeningID,
		})
	}

	c.JSON(200, gin.H{"data": true})
}

// AutoScreenSingle performs auto screening for a single screening record
// POST /api/screening/:id/auto-screen
func AutoScreenSingle(c *gin.Context) {
	id := c.Param("id")

	// Parse screening ID
	var screeningID uint
	if _, err := fmt.Sscanf(id, "%d", &screeningID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid screening ID"})
		return
	}

	// Get admin info from token
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var admin entity.AdminProfile
	if err := config.DB.Where("user_id = ?", userID).First(&admin).Error; err != nil {
		c.JSON(400, gin.H{"error": "Admin profile not found"})
		return
	}

	adminName := admin.AdminFirstname
	if adminName == "" {
		adminName = "ระบบอัตโนมัติ"
	}

	// Perform auto screening
	result, err := services.AutoScreenSingle(screeningID, admin.ID, adminName)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Broadcast to WebSocket for real-time update
	if ws.ApprovalHubInstance != nil {
		log.Printf("📢 Broadcasting auto-screen single completed: ScreeningID=%d", screeningID)
		ws.ApprovalHubInstance.BroadcastUpdate("screening_status_updated", gin.H{
			"screening_id":  screeningID,
			"auto_approved": result.AutoApproved,
		})
	}

	c.JSON(200, gin.H{
		"data":    result,
		"message": getMessage(result.AutoApproved),
	})
}

// AutoScreenBatch performs auto screening for all pending screenings of a scholarship
// POST /api/scholarship/:id/auto-screen-batch
func AutoScreenBatch(c *gin.Context) {
	id := c.Param("id")

	// Parse scholarship ID
	var scholarshipID uint
	if _, err := fmt.Sscanf(id, "%d", &scholarshipID); err != nil {
		c.JSON(400, gin.H{"error": "Invalid scholarship ID"})
		return
	}

	// Check if scholarship exists
	var scholarship entity.Scholarship
	if err := config.DB.First(&scholarship, scholarshipID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Scholarship not found"})
		return
	}

	// Get admin info from token
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var admin entity.AdminProfile
	if err := config.DB.Where("user_id = ?", userID).First(&admin).Error; err != nil {
		c.JSON(400, gin.H{"error": "Admin profile not found"})
		return
	}

	adminName := admin.AdminFirstname
	if adminName == "" {
		adminName = "ระบบอัตโนมัติ"
	}

	// Perform batch auto screening
	results, err := services.AutoScreenBatch(scholarshipID, admin.ID, adminName)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Count results
	autoApproved := 0
	needsReview := 0
	for _, r := range results {
		if r.AutoApproved {
			autoApproved++
		} else {
			needsReview++
		}
	}

	// Broadcast to WebSocket for real-time update
	if ws.ApprovalHubInstance != nil {
		log.Printf("📢 Broadcasting auto-screen batch completed: ScholarshipID=%d, Total=%d", scholarshipID, len(results))
		ws.ApprovalHubInstance.BroadcastUpdate("screening_batch_updated", gin.H{
			"scholarship_id": scholarshipID,
			"total":          len(results),
			"auto_approved":  autoApproved,
			"needs_review":   needsReview,
		})
	}

	c.JSON(200, gin.H{
		"data": results,
		"summary": gin.H{
			"total":         len(results),
			"auto_approved": autoApproved,
			"needs_review":  needsReview,
			"scholarship":   scholarship.ScholarshipName,
		},
		"message": fmt.Sprintf("คัดกรองอัตโนมัติเสร็จสิ้น: อนุมัติ %d คน, รอตรวจสอบ %d คน", autoApproved, needsReview),
	})
}

// getMessage returns appropriate message based on result
func getMessage(autoApproved bool) string {
	if autoApproved {
		return "ผ่านทุกเกณฑ์ - อนุมัติอัตโนมัติแล้ว"
	}
	return "ไม่ผ่านบางเกณฑ์ - รอการตรวจสอบจากเจ้าหน้าที่"
}
