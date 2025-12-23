package controllers

import (
	"backend/config"
	"backend/entity"
	"backend/validators"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Evaluation (การประเมินผู้สมัคร)

// GET /evaluations
func GetAllEvaluations(ctx *gin.Context) {
	var evaluations []entity.Evaluation

	query := config.DB.
		Preload("InterviewRound").
		Preload("ApplicationScholarship.Application.StudentProfile").
		Preload("ApplicationScholarship.Scholarship").
		Preload("AdminProfile").
		Preload("EvaluationScores.EvaluationCriterion")

	// Filter by interview_round_id if provided
	if roundID := ctx.Query("interview_round_id"); roundID != "" {
		query = query.Where("interview_round_id = ?", roundID)
	}

	// Filter by application_scholarship_id if provided
	if appScholarshipID := ctx.Query("application_scholarship_id"); appScholarshipID != "" {
		query = query.Where("application_scholarship_id = ?", appScholarshipID)
	}

	// Filter by status if provided
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// SELECT ALL
	if err := query.Find(&evaluations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, evaluations)
}

// GET /evaluations/:id
func GetEvaluationByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid evaluation id"})
		return
	}

	var evaluation entity.Evaluation

	// SELECT FIRST ROW
	if err := config.DB.
		Preload("InterviewRound").
		Preload("ApplicationScholarship.Application.StudentProfile").
		Preload("ApplicationScholarship.Scholarship").
		Preload("AdminProfile").
		Preload("EvaluationScores.EvaluationCriterion").
		Preload("EvaluationScores.ScoringAdmin").
		First(&evaluation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "evaluation not found"})
		return
	}

	ctx.JSON(http.StatusOK, evaluation)
}

// POST /evaluations
func CreateEvaluation(ctx *gin.Context) {
	var inputValues struct {
		InterviewRoundID         uint   `json:"interview_round_id"`
		ApplicationScholarshipID uint   `json:"application_scholarship_id"`
		AdminID                  uint   `json:"admin_id"`
		Remark                   string `json:"remark"`
	}

	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	evaluation := entity.Evaluation{
		InterviewRoundID:         inputValues.InterviewRoundID,
		ApplicationScholarshipID: inputValues.ApplicationScholarshipID,
		AdminID:                  inputValues.AdminID,
		Status:                   entity.EvaluationStatusPending,
		TotalScore:               0,
		Remark:                   inputValues.Remark,
	}

	// Validate
	if err := validators.ValidateStruct(&evaluation); err != nil {
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
	if err := tx.Create(&evaluation).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.
		Preload("InterviewRound").
		Preload("ApplicationScholarship").
		Preload("AdminProfile").
		First(&evaluation, evaluation.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload"})
		return
	}

	ctx.JSON(http.StatusCreated, evaluation)
}

// PUT /evaluations/:id
func UpdateEvaluation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid evaluation id"})
		return
	}

	var evaluation entity.Evaluation
	if err := config.DB.First(&evaluation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "evaluation not found"})
		return
	}

	var inputValues struct {
		Status *string `json:"status"`
		Remark *string `json:"remark"`
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
	if inputValues.Status != nil {
		updates["status"] = *inputValues.Status
	}
	if inputValues.Remark != nil {
		updates["remark"] = *inputValues.Remark
	}

	if len(updates) > 0 {
		// UPDATE
		if err := tx.Model(&evaluation).Updates(updates).Error; err != nil {
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

	if err := config.DB.
		Preload("InterviewRound").
		Preload("ApplicationScholarship").
		Preload("AdminProfile").
		Preload("EvaluationScores").
		First(&evaluation, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "reload failed"})
		return
	}

	ctx.JSON(http.StatusOK, evaluation)
}

// DELETE /evaluations/:id
func DeleteEvaluation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid evaluation id"})
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

	var evaluation entity.Evaluation
	if err := tx.First(&evaluation, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "evaluation not found"})
		return
	}

	// Delete related scores first
	if err := tx.Where("evaluation_id = ?", id).Delete(&entity.EvaluationScore{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// DELETE
	if err := tx.Delete(&entity.Evaluation{}, id).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Evaluation deleted successfully"})
}

// POST /evaluations/:id/complete
func CompleteEvaluation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid evaluation id"})
		return
	}

	var inputValues struct {
		FinalDecision string `json:"final_decision"` // approved, rejected, waitlist
	}

	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var evaluation entity.Evaluation
	if err := config.DB.First(&evaluation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "evaluation not found"})
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

	// Update evaluation status
	var newStatus entity.EvaluationStatus
	switch inputValues.FinalDecision {
	case "approved":
		newStatus = entity.EvaluationStatusApproved
	case "rejected":
		newStatus = entity.EvaluationStatusRejected
	default:
		newStatus = entity.EvaluationStatusCompleted
	}

	if err := tx.Model(&evaluation).Update("status", newStatus).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update ApplicationScholarship final decision and total score
	if err := tx.Model(&entity.ApplicationScholarship{}).
		Where("id = ?", evaluation.ApplicationScholarshipID).
		Updates(map[string]interface{}{
			"final_total_score": evaluation.TotalScore,
			"final_decision":    inputValues.FinalDecision,
		}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.
		Preload("InterviewRound").
		Preload("ApplicationScholarship").
		Preload("AdminProfile").
		Preload("EvaluationScores").
		First(&evaluation, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "reload failed"})
		return
	}

	ctx.JSON(http.StatusOK, evaluation)
}
