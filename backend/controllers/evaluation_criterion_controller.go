package controllers

import (
	"backend/config"
	"backend/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Evaluation Criterion (เกณฑ์การประเมิน)

// GET /evaluation-criteria
func GetAllEvaluationCriteria(ctx *gin.Context) {
	var criteria []entity.EvaluationCriterion

	// SELECT ALL
	if err := config.DB.Find(&criteria).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

// GET /evaluation-criteria/:id
func GetEvaluationCriterionByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid criterion id"})
		return
	}

	var criterion entity.EvaluationCriterion

	// SELECT FIRST ROW
	if err := config.DB.First(&criterion, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "criterion not found"})
		return
	}

	ctx.JSON(http.StatusOK, criterion)
}

// POST /evaluation-criteria
func CreateEvaluationCriterion(ctx *gin.Context) {
	var inputValues struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		ScoreType   string  `json:"score_type"`
		MaxScore    float64 `json:"max_score"`
		Weight      float64 `json:"weight"`
		IsActive    bool    `json:"is_active"`
	}

	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default values
	scoreType := entity.ScoreTypeNumeric
	if inputValues.ScoreType != "" {
		scoreType = entity.ScoreType(inputValues.ScoreType)
	}
	maxScore := inputValues.MaxScore
	if maxScore == 0 {
		maxScore = 100
	}
	weight := inputValues.Weight
	if weight == 0 {
		weight = 1.0
	}

	criterion := entity.EvaluationCriterion{
		Name:        inputValues.Name,
		Description: inputValues.Description,
		ScoreType:   scoreType,
		MaxScore:    maxScore,
		Weight:      weight,
		IsActive:    inputValues.IsActive,
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
	if err := tx.Create(&criterion).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, criterion)
}

// PUT /evaluation-criteria/:id
func UpdateEvaluationCriterion(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid criterion id"})
		return
	}

	var criterion entity.EvaluationCriterion
	// หา criterion ว่ามีมั้ย
	if err := config.DB.First(&criterion, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "criterion not found"})
		return
	}

	var inputValues struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		ScoreType   *string  `json:"score_type"`
		MaxScore    *float64 `json:"max_score"`
		Weight      *float64 `json:"weight"`
		IsActive    *bool    `json:"is_active"`
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
	// defer เพื่อจัดการ Rollback หากเกิด Panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	updates := map[string]interface{}{}
	if inputValues.Name != nil {
		updates["name"] = *inputValues.Name
	}
	if inputValues.Description != nil {
		updates["description"] = *inputValues.Description
	}
	if inputValues.ScoreType != nil {
		updates["score_type"] = *inputValues.ScoreType
	}
	if inputValues.MaxScore != nil {
		updates["max_score"] = *inputValues.MaxScore
	}
	if inputValues.Weight != nil {
		updates["weight"] = *inputValues.Weight
	}
	if inputValues.IsActive != nil {
		updates["is_active"] = *inputValues.IsActive
	}

	if len(updates) > 0 {
		// UPDATE
		if err := tx.Model(&criterion).Updates(updates).Error; err != nil {
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

	if err := config.DB.First(&criterion, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "reload failed"})
		return
	}

	ctx.JSON(http.StatusOK, criterion)
}

// DELETE /evaluation-criteria/:id
func DeleteEvaluationCriterion(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid criterion id"})
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

	var criterion entity.EvaluationCriterion
	if err := tx.First(&criterion, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "criterion not found"})
		return
	}

	// DELETE
	if err := tx.Delete(&entity.EvaluationCriterion{}, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Criterion deleted successfully"})
}
