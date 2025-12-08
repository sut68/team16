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

	if err := config.DB.Preload("StudentProfile").Preload("ApplicationDocuments").Find(&applications).Error; err != nil {
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

	if err := config.DB.Preload("StudentProfile").Preload("ApplicationDocuments").First(&application, id).Error; err != nil {
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
