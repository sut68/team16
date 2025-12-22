package controllers

import (
	"backend/config"
	"backend/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Evaluation Score (คะแนนรายเกณฑ์)

// POST /evaluations/:id/scores
func AddEvaluationScore(ctx *gin.Context) {
	evalIDStr := ctx.Param("id")
	evalID, err := strconv.Atoi(evalIDStr)
	if err != nil || evalID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid evaluation id"})
		return
	}

	var inputValues struct {
		EvaluationCriterionID uint    `json:"evaluation_criterion_id"`
		ScoreValue            float64 `json:"score_value"`
		Comment               string  `json:"comment"`
		ScoringAdminID        uint    `json:"scoring_admin_id"`
	}

	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score := entity.EvaluationScore{
		EvaluationID:          uint(evalID),
		EvaluationCriterionID: inputValues.EvaluationCriterionID,
		ScoreValue:            inputValues.ScoreValue,
		Comment:               inputValues.Comment,
		ScoringAdminID:        inputValues.ScoringAdminID,
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

	// INSERT
	if err := tx.Create(&score).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Recalculate total score
	RecalculateTotalScore(uint(evalID))

	if err := config.DB.Preload("EvaluationCriterion").First(&score, score.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload"})
		return
	}

	ctx.JSON(http.StatusCreated, score)
}

// PUT /evaluation-scores/:id
func UpdateEvaluationScore(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid score id"})
		return
	}

	var score entity.EvaluationScore
	if err := config.DB.First(&score, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "score not found"})
		return
	}

	var inputValues struct {
		ScoreValue *float64 `json:"score_value"`
		Comment    *string  `json:"comment"`
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
	if inputValues.ScoreValue != nil {
		updates["score_value"] = *inputValues.ScoreValue
	}
	if inputValues.Comment != nil {
		updates["comment"] = *inputValues.Comment
	}

	if len(updates) > 0 {
		// UPDATE
		if err := tx.Model(&score).Updates(updates).Error; err != nil {
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

	// Recalculate total score
	RecalculateTotalScore(score.EvaluationID)

	if err := config.DB.Preload("EvaluationCriterion").First(&score, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "reload failed"})
		return
	}

	ctx.JSON(http.StatusOK, score)
}

// DELETE /evaluation-scores/:id
func DeleteEvaluationScore(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid score id"})
		return
	}

	var score entity.EvaluationScore
	if err := config.DB.First(&score, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "score not found"})
		return
	}

	evalID := score.EvaluationID

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

	// DELETE
	if err := tx.Delete(&entity.EvaluationScore{}, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Recalculate total score
	RecalculateTotalScore(evalID)

	ctx.JSON(http.StatusOK, gin.H{"message": "Score deleted successfully"})
}

// RecalculateTotalScore - คำนวณคะแนนรวมสำหรับการประเมิน
func RecalculateTotalScore(evaluationID uint) {
	var scores []entity.EvaluationScore
	config.DB.Where("evaluation_id = ?", evaluationID).Preload("EvaluationCriterion").Find(&scores)

	var totalScore float64
	var totalWeight float64

	for _, score := range scores {
		criterion := score.EvaluationCriterion
		// Normalize score to percentage and apply weight
		normalizedScore := (score.ScoreValue / criterion.MaxScore) * 100
		weightedScore := normalizedScore * criterion.Weight
		totalScore += weightedScore
		totalWeight += criterion.Weight
	}

	// Calculate weighted average
	var finalScore float64
	if totalWeight > 0 {
		finalScore = totalScore / totalWeight
	}

	// Update evaluation
	config.DB.Model(&entity.Evaluation{}).Where("id = ?", evaluationID).Update("total_score", finalScore)

	// Update status to in_progress if scores exist
	if len(scores) > 0 {
		config.DB.Model(&entity.Evaluation{}).Where("id = ? AND status = ?", evaluationID, entity.EvaluationStatusPending).Update("status", entity.EvaluationStatusInProgress)
	}
}
