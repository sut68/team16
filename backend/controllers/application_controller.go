package controllers

import (
	"backend/config"
	"backend/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /applications
func GetApplications(ctx *gin.Context) {
	var applications []entity.Application

	if err := config.DB.Preload("StudentProfile").Preload("ApplicationScholarships").Find(&applications).Error; err != nil {
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

	if err := config.DB.Preload("StudentProfile").Preload("ApplicationScholarships").First(&application, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	ctx.JSON(http.StatusOK, application)
}

// POST /applications
func CreateApplication(ctx *gin.Context) {
	var input struct {
		StudentProfileID uint `json:"student_profile_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application := entity.Application{
		StudentProfileID: input.StudentProfileID,
	}

	if err := config.DB.Create(&application).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, application)
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

	// Find the parent Application record for the student
	var application entity.Application
	if err := config.DB.Where("student_profile_id = ?", studentID).First(&application).Error; err != nil {
		// If the student has never applied for anything, they won't have an Application record.
		// Return an empty array instead of a 404.
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

	ctx.JSON(http.StatusOK, appScholarships)
}
