package controllers

import (
	"backend/config"
	"backend/entity"
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getall
func GetAllScholarship(c *gin.Context) {
	// Define a custom response struct to avoid circular dependencies during JSON marshaling.
	// This struct should only contain fields needed for the list view.
	type ScholarshipResponse struct {
		ID                uint                     `json:"ID"`
		ScholarshipName   string                   `json:"scholarship_name"`
		Description       string                   `json:"description"`
		OpenDate          string                   `json:"open_date"`
		CloseDate         string                   `json:"close_date"`
		Statusscholarship entity.Statusscholarship `json:"statusscholarship"`
		Typescholarship   entity.Typescholarship   `json:"typescholarship"`
		Semaster          entity.Semaster          `json:"semaster"`
		Sponsor           entity.Sponsor           `json:"sponsor"`
	}

	var scholarships []entity.Scholarship
	// Only preload the data that is actually needed for the list view.
	// Preloading ApprovalRequirements or ApplicationScholarships here would cause a crash.
	if err := config.DB.
		Preload("Statusscholarship").
		Preload("Typescholarship").
		Preload("Semaster").
		Preload("Sponsor").
		Find(&scholarships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Manually map to the clean response struct to ensure no circular references are passed to the JSON marshaller.
	var response []ScholarshipResponse
	for _, s := range scholarships {
		response = append(response, ScholarshipResponse{
			ID:                s.ID,
			ScholarshipName:   s.ScholarshipName,
			Description:       s.Description,
			OpenDate:          s.OpenDate,
			CloseDate:         s.CloseDate,
			Statusscholarship: s.Statusscholarship,
			Typescholarship:   s.Typescholarship,
			Semaster:          s.Semaster,
			Sponsor:           s.Sponsor,
		})
	}

	c.JSON(http.StatusOK, response)
}

// getby id
func GetScholarshipByID(c *gin.Context) {
	id := c.Param("id")
	var item entity.Scholarship
	if err := config.DB.Preload("Statusscholarship").Preload("Typescholarship").Preload("Semaster").Preload("Sponsor").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// post
func CreateScholarship(c *gin.Context) {
	var item entity.Scholarship
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// put
func UpdateScholarship(c *gin.Context) {
	id := c.Param("id")
	var item entity.Scholarship
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scholarship not found"})
		return
	}

	var input entity.Scholarship
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&item).Updates(input)
	c.JSON(http.StatusOK, item)
}

// delete
func DeleteScholarship(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&entity.Scholarship{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Scholarship deleted"})
}

// POST /scholarships/:id/apply
func ApplyForScholarship(ctx *gin.Context) {
	// Start a new database transaction
	tx := config.DB.Begin()

	// --- 1. Get Input ---
	scholarshipIDStr := ctx.Param("id")
	scholarshipID, err := strconv.ParseUint(scholarshipIDStr, 10, 64)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid scholarship id"})
		return
	}

	// Extended input struct to receive more data
	var input struct {
		StudentProfileID  uint   `json:"student_profile_id" binding:"required"`
		ApplicationReason string `json:"application_reason"` // เหตุผลในการสมัครทุน

		// ข้อมูลที่อนุญาตให้แก้ไข
		Email string `json:"email"`
		Phone string `json:"phone"`

		// ข้อมูลครอบครัว
		FatherOccupation   string  `json:"father_occupation"`
		FatherIncome       float64 `json:"father_income"`
		MotherOccupation   string  `json:"mother_occupation"`
		MotherIncome       float64 `json:"mother_income"`
		GuardianOccupation string  `json:"guardian_occupation"`
		GuardianIncome     float64 `json:"guardian_income"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- 2. Update student profile if email/phone provided ---
	if input.Email != "" || input.Phone != "" {
		var student entity.StudentProfile
		if err := tx.First(&student, input.StudentProfileID).Error; err == nil {
			if input.Email != "" {
				student.Email = input.Email
			}
			if input.Phone != "" {
				student.Phone = input.Phone
			}
			tx.Save(&student)
		}
	}

	// --- 3. Update family info if provided ---
	if input.FatherOccupation != "" || input.MotherOccupation != "" ||
		input.FatherIncome > 0 || input.MotherIncome > 0 {
		var family entity.FamilyInfo
		err := tx.Where("profile_id = ?", input.StudentProfileID).First(&family).Error
		if err == nil {
			// Update existing family info
			if input.FatherOccupation != "" {
				family.FatherOccupation = input.FatherOccupation
			}
			if input.FatherIncome > 0 {
				family.FatherIncome = input.FatherIncome
			}
			if input.MotherOccupation != "" {
				family.MotherOccupation = input.MotherOccupation
			}
			if input.MotherIncome > 0 {
				family.MotherIncome = input.MotherIncome
			}
			if input.GuardianOccupation != "" {
				family.GuardianOccupation = input.GuardianOccupation
			}
			if input.GuardianIncome > 0 {
				family.GuardianIncome = input.GuardianIncome
			}
			tx.Save(&family)
		}
	}

	// --- 4. Find the currently active semester ---
	var activeSemester entity.Semaster
	if err := tx.Where("is_active = ?", true).First(&activeSemester).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No active semester found. Cannot process application."})
		return
	}

	// --- 5. Find or Create parent Application record for the student in the active semester ---
	application := entity.Application{
		StudentProfileID: input.StudentProfileID,
		SemasterID:       activeSemester.ID,
	}
	if err := tx.Where(entity.Application{
		StudentProfileID: input.StudentProfileID,
		SemasterID:       activeSemester.ID,
	}).FirstOrCreate(&application).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find or create application: " + err.Error()})
		return
	}

	// --- 6. Create the ApplicationScholarship join record with reason ---
	var reasonPtr *string
	if input.ApplicationReason != "" {
		reasonPtr = &input.ApplicationReason
	}

	appScholarship := entity.ApplicationScholarship{
		ApplicationID:     application.ID,
		ScholarshipID:     uint(scholarshipID),
		Status:            "new", // 'new' status signifies that the application is awaiting screening.
		ApplicationReason: reasonPtr,
	}
	if err := tx.Create(&appScholarship).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application_scholarship record: " + err.Error()})
		return
	}

	// --- 7. Create Screening record for admin to review ---
	screening := entity.Screening{
		AdminProfileID:           1, // Default to first admin (will be reassigned when admin reviews)
		ApplicationScholarshipID: appScholarship.ID,
		StatusScreeningID:        1, // 1 = "รอตรวจสอบ" (Pending Review)
	}
	if err := tx.Create(&screening).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create screening record: " + err.Error()})
		return
	}

	// --- 8. Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed: " + err.Error()})
		return
	}

	// --- 9. Auto Screening (async) ---
	// Run auto screening in background to not block the response
	// ผู้สมัครจะได้รับการคัดกรองทันที ถ้าผ่านทุกเกณฑ์จะ auto approve
	go func(screeningID uint) {
		result, err := services.AutoScreenSingle(screeningID, 1, "ระบบอัตโนมัติ")
		if err != nil {
			// Log error but don't fail - admin can review manually
			// Note: This uses log package, make sure it's imported
			return
		}

		// Log the result
		if result.AutoApproved {
			// Auto approved - student passed all criteria
			_ = result // Result is broadcast via WebSocket in AutoScreenSingle
		}
		// If not auto approved, screening stays pending for manual review
	}(screening.ID)

	ctx.JSON(http.StatusCreated, gin.H{
		"message":                  "Application created successfully. Auto screening in progress.",
		"applicationId":            application.ID,
		"applicationScholarshipId": appScholarship.ID,
		"screeningId":              screening.ID,
	})
}

// GetQualifiedApplicantsForScholarship godoc
// GET /scholarships/:id/qualified-applicants
func GetQualifiedApplicantsForScholarship(c *gin.Context) {
	scholarshipIDStr := c.Param("id")
	scholarshipID, err := strconv.Atoi(scholarshipIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scholarship ID"})
		return
	}

	var qualifiedApplicants []entity.ApplicationScholarship

	// Subquery to find all application_scholarship_ids that already have a booking
	subQuery := config.DB.Model(&entity.IntervieweBooking{}).Select("application_scholarship_id")

	if err := config.DB.
		Where("scholarship_id = ?", scholarshipID).
		Where("status = ?", "qualified").
		Where("id NOT IN (?)", subQuery).
		Preload("Application.StudentProfile.Major").
		Find(&qualifiedApplicants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query qualified applicants: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, qualifiedApplicants)
}
