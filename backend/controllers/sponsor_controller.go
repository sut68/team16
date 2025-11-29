package controllers

import (
	"backend/config"
	"backend/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /sponsors
func GetSponsors(ctx *gin.Context) {
	var sponsors []entity.Sponsor

	// SELECT ALL
	if err := config.DB.Preload("Industry").Preload("Contacts").Find(&sponsors).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sponsors)
}

// GET /sponsors/:id 
func GetSponsorsByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid sponsor id"})
		return
	}

	var sponsor entity.Sponsor

	// SELECT FIRST ROW
	if err := config.DB.Preload("Industry").Preload("Contacts").First(&sponsor, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sponsor)
}

// POST /sponsors
func CreateSponsor(ctx *gin.Context) {
	var inputValues struct {
		CompanyName		string		`json:"company_name" binding:"required"`
		IndustryID		*uint     `json:"industry_id"`
		Website				*string		`json:"website"`
		Status				string		`json:"status"`
		Description		*string		`json:"description"`
		Contacts			[]entity.SponsorContact		`json:"contacts"`
	}

	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sponsor := entity.Sponsor{
		CompanyName:		inputValues.CompanyName,
		IndustryID:			inputValues.IndustryID,
		Website:				inputValues.Website,
		Status:					inputValues.Status,
		Description:		inputValues.Description,
		Contacts:			inputValues.Contacts,
	}

	// INSERT
	if err := config.DB.Create(&sponsor).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, sponsor)
}

// PATCH /sponsors/:id
func UpdateSponsor(ctx *gin.Context) {
	id := ctx.Param("id")
	var sponsor entity.Sponsor

	// หา sponsor ว่ามีมั้ย กันส่ง มาไม่มีอยู่
	if err := config.DB.Preload("Contacts").First(&sponsor, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!!!!"})
		return
	}

	var inputValues struct {
		CompanyName		*string		`json:"company_name"`
		IndustryID    *string		`json:"industry_id"`
		Website				*string		`json:"website"`
		Status				*string		`json:"status"`
		Description		*string		`json:"description"`
		Contacts			[]entity.SponsorContact		`json:"contacts"`
	}

	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if inputValues.CompanyName != nil {
		updates["company_name"] = *inputValues.CompanyName
	}
	if inputValues.IndustryID != nil {
		updates["industry_id"] = *inputValues.IndustryID
	}
	if inputValues.Website != nil {
		updates["website"] = *inputValues.Website
	}
	if inputValues.Status != nil {
		updates["status"] = *inputValues.Status
	}
	if inputValues.Description != nil {
		updates["description"] = *inputValues.Description
	}

	if len(updates) > 0 {
		// UPDATE
		if err := config.DB.Model(&sponsor).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Contact
	if inputValues.Contacts != nil {
		config.DB.Where("sponsor_id = ?", sponsor.ID).Delete(&entity.SponsorContact{})
		for i := range inputValues.Contacts {
			inputValues.Contacts[i].SponsorID = sponsor.ID
		}

		if err := config.DB.Create(&inputValues.Contacts).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, sponsor)
}

// DELETE /sponsors/:id
func DeleteSponsor(ctx *gin.Context) {
	id := ctx.Param("id")

	// DELETE
	if err := config.DB.Delete(&entity.Sponsor{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sponsor delete successfully"})
}