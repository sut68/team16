package controllers

import (
	"backend/config"
	"backend/entity"
	"github.com/gin-gonic/gin"
)

func GetAllScreenings(c *gin.Context) {
	var screenings []entity.Screening

	err := config.DB.
		Preload("StatusScreening").
		Preload("ApplicationScholarship.Application.StudentProfile").
		Preload("ApplicationScholarship.Scholarship").
        Preload("ApplicationScholarship.Scholarship.Semaster").
		Find(&screenings).Error

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": screenings})
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
        Preload("ApplicationScholarship.Scholarship"). // ไม่โหลด Sponsor
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
            "ID":                   screening.ID,
            "CreatedAt":            screening.CreatedAt,
            "UpdatedAt":            screening.UpdatedAt,
            "DeletedAt":            screening.DeletedAt,
            "admin_profile_id":     screening.AdminProfileID,
            "admin_profile":        screening.AdminProfile,
            "application_scholarship_id": screening.ApplicationScholarshipID,
            "application_scholarship": gin.H{
                "ID":           screening.ApplicationScholarship.ID,
                "status":       screening.ApplicationScholarship.Status,
                "application_id": screening.ApplicationScholarship.ApplicationID,
                "application": screening.ApplicationScholarship.Application,
                "scholarship_id": screening.ApplicationScholarship.ScholarshipID,
                "scholarship": gin.H{
                    "ID":                 screening.ApplicationScholarship.Scholarship.ID,
                    "scholarship_name":   screening.ApplicationScholarship.Scholarship.ScholarshipName,
                    "description":        screening.ApplicationScholarship.Scholarship.Description,
                    "open_date":          screening.ApplicationScholarship.Scholarship.OpenDate,
                    "close_date":         screening.ApplicationScholarship.Scholarship.CloseDate,
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

	const (
		PASS = 1
		FAIL = 2
	)

	if input.StatusScreeningID == FAIL && input.RejectionReason == nil {
		c.JSON(400, gin.H{"error": "Rejection reason is required"})
		return
	}

	if input.StatusScreeningID == PASS {
		input.RejectionReason = nil
	}

	if err := config.DB.Model(&screening).Updates(map[string]interface{}{
		"status_screening_id": input.StatusScreeningID,
		"rejection_reason":    input.RejectionReason,
	}).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": true})
}
