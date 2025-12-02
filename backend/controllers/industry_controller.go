package controllers

import (
	"backend/config"
	"backend/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /industries
func GetIndustries(ctx *gin.Context) {
	var industries []entity.SponsorIndustry
	if err := config.DB.Find(&industries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	
	ctx.JSON(http.StatusOK, industries)
}

// POST /industries
func CreateIndustry(ctx *gin.Context) {
	var inputValues entity.SponsorIndustry
	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&inputValues).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, inputValues)
}