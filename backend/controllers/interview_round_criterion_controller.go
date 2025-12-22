package controllers

import (
	"backend/config"
	"backend/entity"
	"backend/validators"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Interview Round Criterion (เกณฑ์การประเมินรอบสัมภาษณ์)

// GET /interview-rounds/:id/criteria
func GetInterviewRoundCriteria(ctx *gin.Context) {
	roundIDStr := ctx.Param("id")
	roundID, err := strconv.Atoi(roundIDStr)
	if err != nil || roundID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid interview round id"})
		return
	}

	var criteria []entity.InterviewRoundCriterion

	// SELECT ALL with Preload
	if err := config.DB.
		Where("interview_round_id = ?", roundID).
		Preload("EvaluationCriterion").
		Find(&criteria).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

// POST /interview-rounds/:id/criteria
func AddCriterionToInterviewRound(ctx *gin.Context) {
	roundIDStr := ctx.Param("id")
	roundID, err := strconv.Atoi(roundIDStr)
	if err != nil || roundID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid interview round id"})
		return
	}

	var inputValues struct {
		EvaluationCriterionID uint    `json:"evaluation_criterion_id"`
		Weight                float64 `json:"weight"`
		IsEnabled             bool    `json:"is_enabled"`
	}

	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default values
	weight := inputValues.Weight
	if weight == 0 {
		weight = 1.0
	}

	roundCriterion := entity.InterviewRoundCriterion{
		InterviewRoundID:      uint(roundID),
		EvaluationCriterionID: inputValues.EvaluationCriterionID,
		Weight:                weight,
		IsEnabled:             inputValues.IsEnabled,
	}

	// Validate
	if err := validators.ValidateStruct(&roundCriterion); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"error":   err.Error(),
		})
		return
	}

	// เริ่มต้นการดำเนินการฐานข้อมูล
	tx := config.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	// defer เพื่อจัดการ Rollback หากเกิด Panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// INSERT
	if err := tx.Create(&roundCriterion).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Preload("EvaluationCriterion").First(&roundCriterion, roundCriterion.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload"})
		return
	}

	ctx.JSON(http.StatusCreated, roundCriterion)
}

// PUT /interview-round-criteria/:id
func UpdateInterviewRoundCriterion(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	var roundCriterion entity.InterviewRoundCriterion
	if err := config.DB.First(&roundCriterion, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "round criterion not found"})
		return
	}

	var inputValues struct {
		Weight    *float64 `json:"weight"`
		IsEnabled *bool    `json:"is_enabled"`
	}

	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// เริ่มต้นการดำเนินการฐานข้อมูล
	tx := config.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	updates := map[string]interface{}{}
	if inputValues.Weight != nil {
		updates["weight"] = *inputValues.Weight
	}
	if inputValues.IsEnabled != nil {
		updates["is_enabled"] = *inputValues.IsEnabled
	}

	if len(updates) > 0 {
		// UPDATE
		if err := tx.Model(&roundCriterion).Updates(updates).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Preload("EvaluationCriterion").First(&roundCriterion, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "reload failed"})
		return
	}

	ctx.JSON(http.StatusOK, roundCriterion)
}

// DELETE /interview-round-criteria/:id
func RemoveCriterionFromInterviewRound(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var roundCriterion entity.InterviewRoundCriterion
	if err := tx.First(&roundCriterion, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "round criterion not found"})
		return
	}

	// DELETE
	if err := tx.Delete(&entity.InterviewRoundCriterion{}, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Criterion removed from interview round"})
}
