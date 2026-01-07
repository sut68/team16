package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/entity"
)

// GET /applications
func GetApplications(ctx *gin.Context) {
	var applications []entity.Application

	if err := config.DB.Preload("StudentProfile").Preload("ApplicationScholarships").Preload("Semaster").Find(&applications).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, applications)
}

// GET /application-scholarships
// สามารถ filter ด้วย status ได้ เช่น ?status=qualified
func GetAllApplicationScholarships(ctx *gin.Context) {
	var appScholarships []entity.ApplicationScholarship

	query := config.DB.
		Preload("Application.StudentProfile.Major").
		Preload("Scholarship.Statusscholarship").
		Preload("Scholarship.Typescholarship").
		Preload("ApplicationDocuments.ApprovalTasks.ApprovalDecisions").
		Preload("Screening.StatusScreening").
		Preload("IntervieweBookings.Slot").
		Preload("Evaluations")

	// Filter by status
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&appScholarships).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, appScholarships)
}

// GET /applications/:id
func GetApplicationByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid application id"})
		return
	}

	var application entity.Application

	if err := config.DB.Preload("StudentProfile").Preload("ApplicationScholarships").Preload("Semaster").First(&application, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	ctx.JSON(http.StatusOK, application)
}

// POST /applications
func CreateApplication(ctx *gin.Context) {
	var input struct {
		StudentProfileID uint `json:"student_profile_id" binding:"required"`
		ScholarshipID    uint `json:"scholarship_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	// Find the active semester
	var activeSemester entity.Semaster
	if err := tx.Where("is_active = ?", true).First(&activeSemester).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No active semester found"})
		return
	}

	// Find or create the main application record for the student and semester
	var application entity.Application
	err := tx.Where("student_profile_id = ? AND semaster_id = ?", input.StudentProfileID, activeSemester.ID).
		FirstOrCreate(&application, entity.Application{
			StudentProfileID: input.StudentProfileID,
			SemasterID:       activeSemester.ID,
		}).Error

	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find or create application: " + err.Error()})
		return
	}

	// Create the specific scholarship application entry
	appScholarship := entity.ApplicationScholarship{
		ApplicationID: application.ID,
		ScholarshipID: input.ScholarshipID,
	}

	if err := tx.Create(&appScholarship).Error; err != nil {
		tx.Rollback()
		// Check for duplicate entry error
		if IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "You have already applied for this scholarship."})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create scholarship application: " + err.Error()})
		}
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":                 "Application successful!",
		"application":             application,
		"application_scholarship": appScholarship,
	})
}

// DELETE /applications/:id
func DeleteApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid application id"})
		return
	}

	if err := config.DB.Delete(&entity.Application{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
}

// GET /students/:student_profile_id/applications
func GetStudentApplications(ctx *gin.Context) {
	studentIDStr := ctx.Param("student_profile_id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid student profile id"})
		return
	}

	// Find the active semester
	var activeSemester entity.Semaster
	if err := config.DB.Where("is_active = ?", true).First(&activeSemester).Error; err != nil {
		// If no active semester, the student has no active applications
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusOK, []entity.ApplicationScholarship{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find active semester: " + err.Error()})
		return
	}

	// Find the parent Application record for the student FOR THE CURRENT SEMESTER
	var application entity.Application
	if err := config.DB.Where("student_profile_id = ? AND semaster_id = ?", studentID, activeSemester.ID).First(&application).Error; err != nil {
		// If no application for this student in this semester, return empty list
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusOK, []entity.ApplicationScholarship{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find application record: " + err.Error()})
		return
	}

	// Preload all associated scholarship applications
	var appScholarships []entity.ApplicationScholarship
	if err := config.DB.
		Where("application_id = ?", application.ID).
		Preload("Scholarship.Statusscholarship").
		Preload("Scholarship.Typescholarship").
		Preload("Application").
		Preload("ApplicationDocuments.ApprovalTasks.ApprovalDecisions").
		Preload("Screening").
		Preload("Screening.StatusScreening").
		Find(&appScholarships).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scholarship applications: " + err.Error()})
		return
	}

	// Process status for each application
	// Only compute status dynamically if it's not already set in database
	for i := range appScholarships {
		existingStatus := appScholarships[i].Status

		// If status is already determined (qualified, rejected, approved, etc.), keep it
		if existingStatus == "qualified" || existingStatus == "rejected" ||
			existingStatus == "approved" || existingStatus == "completed" {
			continue // Use the status from database
		}

		// For other statuses (new, pending, empty), compute based on documents
		isQualified := true
		if len(appScholarships[i].ApplicationDocuments) == 0 {
			isQualified = false
		}
		for _, doc := range appScholarships[i].ApplicationDocuments {
			isDocApproved := false
			hasReject := false
			for _, task := range doc.ApprovalTasks {
				// Check task status instead of iterating all decisions
				if task.Status == "approved" {
					isDocApproved = true
				}
				if task.Status == "rejected" {
					hasReject = true
				}
			}
			if !isDocApproved || hasReject {
				isQualified = false
				break
			}
		}
		if isQualified && len(appScholarships[i].ApplicationDocuments) > 0 {
			appScholarships[i].Status = "qualified"
		} else if existingStatus == "" {
			appScholarships[i].Status = "pending"
		}
	}

	ctx.JSON(http.StatusOK, appScholarships)
}

func IsDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint failed")
}

// DELETE /application-scholarships/:id/cancel
// ยกเลิกการสมัครทุน - สามารถยกเลิกได้เฉพาะสถานะ new, pending, หรือยังไม่ได้รับการตรวจสอบ
func CancelApplicationScholarship(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application scholarship ID"})
		return
	}

	// Start transaction
	tx := config.DB.Begin()

	// 1. Find the application scholarship
	var appScholarship entity.ApplicationScholarship
	if err := tx.Preload("Screening").First(&appScholarship, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Application scholarship not found"})
		return
	}

	// 2. Check if cancellation is allowed based on status
	allowedStatuses := []string{"new", "pending", ""}
	isAllowed := false
	for _, status := range allowedStatuses {
		if appScholarship.Status == status {
			isAllowed = true
			break
		}
	}

	// Also check screening status - if screening is approved or rejected, cannot cancel
	if appScholarship.Screening != nil {
		screeningStatusId := appScholarship.Screening.StatusScreeningID
		// 1 = pending, 2 = approved, 3 = rejected
		if screeningStatusId == 2 || screeningStatusId == 3 {
			isAllowed = false
		}
	}

	if !isAllowed {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ไม่สามารถยกเลิกได้ เนื่องจากใบสมัครอยู่ในขั้นตอนการพิจารณาหรือดำเนินการแล้ว",
		})
		return
	}

	// 3. Delete related Screening record first (if exists)
	if appScholarship.Screening != nil {
		if err := tx.Delete(&entity.Screening{}, "application_scholarship_id = ?", id).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete screening record: " + err.Error()})
			return
		}
	}

	// 4. Delete related Application Documents (if any)
	if err := tx.Delete(&entity.ApplicationDocument{}, "application_scholarship_id = ?", id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete documents: " + err.Error()})
		return
	}

	// 5. Delete the Application Scholarship record
	if err := tx.Delete(&appScholarship).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel application: " + err.Error()})
		return
	}

	// 6. Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ยกเลิกการสมัครทุนเรียบร้อยแล้ว",
	})
}
