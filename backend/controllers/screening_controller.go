package controllers

import (
	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/entity"
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

	c.JSON(200, gin.H{"data": true})
}
