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
		"message":                   "Application successful!",
		"application":               application,
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
		Find(&appScholarships).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scholarship applications: " + err.Error()})
		return
	}

	for i := range appScholarships {
		isQualified := true
		if len(appScholarships[i].ApplicationDocuments) == 0 {
			isQualified = false
		}
		for _, doc := range appScholarships[i].ApplicationDocuments {
			isDocApproved := false
			hasReject := false
			for _, task := range doc.ApprovalTasks {
				for _, decision := range task.ApprovalDecisions {
					if decision.Decision == "approve" {
						isDocApproved = true
					}
					if decision.Decision == "reject" {
						hasReject = true
					}
				}
			}
			if !isDocApproved || hasReject {
				isQualified = false
				break
			}
		}
		if isQualified {
			appScholarships[i].Status = "qualified"
		} else {
			appScholarships[i].Status = "pending"
		}
	}

	ctx.JSON(http.StatusOK, appScholarships)
}

func IsDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint failed")
}
