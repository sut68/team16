package controllers

import (
	"backend/config"
	"backend/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// --- ApprovalTask Functions ---

// GET /approval-tasks
func GetApprovalTasks(ctx *gin.Context) {
	var tasks []entity.ApprovalTask
	// เพิ่ม Preload แบบซ้อนกัน (Nested Preload) เพื่อดึงข้อมูลที่สัมพันธ์กันออกมาให้ครบ
	if err := config.DB.
		Preload("Admin").
		Preload("ApplicationDocument").
		Preload("Application.StudentProfile.Major"). // ดึงไปถึงสาขาวิชาของนิสิต
		Preload("Application.StudentProfile").       // ดึงข้อมูลนิสิต
		Preload("ApprovalRequirement.Scholarship").  // ดึงข้อมูลชื่อทุน
		Preload("ApprovalRequirement").
		Preload("ApprovalDecisions"). // เพิ่ม Preload ApprovalDecisions
		Find(&tasks).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

// GET /approval-tasks/:id
func GetApprovalTaskByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid approval task id"})
		return
	}
	var task entity.ApprovalTask
	if err := config.DB.
		Preload("Admin").
		Preload("ApplicationDocument").
		Preload("Application.StudentProfile.Major").
		Preload("Application.StudentProfile").
		Preload("ApprovalRequirement.Scholarship").
		Preload("ApprovalRequirement").
		Preload("ApprovalDecisions"). // เพิ่ม Preload ApprovalDecisions
		First(&task, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Approval task not found"})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

// POST /approval-tasks
func CreateApprovalTask(ctx *gin.Context) {
	var input struct {
		Status        string `json:"status" binding:"required"`
		AdminID       uint   `json:"admin_id" binding:"required"`
		DocumentID    uint   `json:"document_id" binding:"required"`
		ApplicationID uint   `json:"application_id" binding:"required"`
		RequirementID uint   `json:"requirement_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := entity.ApprovalTask{
		Status:        input.Status,
		AdminID:       input.AdminID,
		DocumentID:    input.DocumentID,
		ApplicationID: input.ApplicationID,
		RequirementID: input.RequirementID,
	}
	if err := config.DB.Create(&task).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, task)
}

// PATCH /approval-tasks/:id
func UpdateApprovalTask(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid approval task id"})
		return
	}
	var task entity.ApprovalTask
	if err := config.DB.First(&task, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Approval task not found"})
		return
	}
	var input struct {
		Status  *string `json:"status"`
		AdminID *uint   `json:"admin_id"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := make(map[string]interface{})
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.AdminID != nil {
		updates["admin_id"] = *input.AdminID
	}
	if len(updates) > 0 {
		if err := config.DB.Model(&task).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, task)
}

// DELETE /approval-tasks/:id
func DeleteApprovalTask(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid approval task id"})
		return
	}
	if err := config.DB.Delete(&entity.ApprovalTask{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Approval task deleted successfully"})
}

// GET /students/:id/decision-history
func GetDecisionHistoryByStudentID(ctx *gin.Context) {
	studentID := ctx.Param("id")
	if studentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Student ID is required"})
		return
	}
	var decisions []entity.ApprovalDecision
	err := config.DB.
		Joins("JOIN approval_tasks ON approval_tasks.id = approval_decisions.task_id").
		Joins("JOIN application_documents ON application_documents.id = approval_tasks.document_id").
		Joins("JOIN applications ON applications.id = application_documents.application_id").
		Where("applications.student_profile_id = ?", studentID).
		Preload("ApprovalTask").
		Find(&decisions).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(decisions) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No decision history found for this student"})
		return
	}
	ctx.JSON(http.StatusOK, decisions)
}

// --- ApplicationDocument Functions ---

// GET /application-documents
func GetApplicationDocuments(ctx *gin.Context) {
	var documents []entity.ApplicationDocument
	if err := config.DB.Preload("Application").Preload("ApprovalRequirement").Find(&documents).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, documents)
}

// GET /application-documents/:id
func GetApplicationDocumentByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid application document id"})
		return
	}
	var document entity.ApplicationDocument
	if err := config.DB.Preload("Application").Preload("ApprovalRequirement").First(&document, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Application document not found"})
		return
	}
	ctx.JSON(http.StatusOK, document)
}

// POST /application-documents
func CreateApplicationDocument(ctx *gin.Context) {
	var input struct {
		FileName      string `json:"file_name" binding:"required"`
		UploadedBy    string `json:"uploaded_by" binding:"required"`
		ApplicationID uint   `json:"application_id" binding:"required"`
		RequirementID uint   `json:"requirement_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	document := entity.ApplicationDocument{
		FileName:      input.FileName,
		UploadedBy:    input.UploadedBy,
		ApplicationID: input.ApplicationID,
		RequirementID: input.RequirementID,
	}
	if err := config.DB.Create(&document).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, document)
}

// PATCH /application-documents/:id
func UpdateApplicationDocument(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid application document id"})
		return
	}
	var document entity.ApplicationDocument
	if err := config.DB.First(&document, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Application document not found"})
		return
	}
	var input struct {
		FileName *string `json:"file_name"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.FileName != nil {
		if err := config.DB.Model(&document).Update("file_name", *input.FileName).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, document)
}

// DELETE /application-documents/:id
func DeleteApplicationDocument(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid application document id"})
		return
	}
	if err := config.DB.Delete(&entity.ApplicationDocument{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Application document deleted successfully"})
}

// --- ApprovalDecision Functions ---

// GET /approval-decisions
func GetApprovalDecisions(ctx *gin.Context) {
	var decisions []entity.ApprovalDecision
	if err := config.DB.Preload("ApprovalTask").Find(&decisions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, decisions)
}

// GET /approval-decisions/:id
func GetApprovalDecisionByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid approval decision id"})
		return
	}
	var decision entity.ApprovalDecision
	if err := config.DB.Preload("ApprovalTask").First(&decision, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Approval decision not found"})
		return
	}
	ctx.JSON(http.StatusOK, decision)
}

// POST /approval-decisions
func CreateApprovalDecision(ctx *gin.Context) {
	var input struct {
		Decision string `json:"decision" binding:"required"`
		Comment  string `json:"comment"`
		TaskID   uint   `json:"task_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	decision := entity.ApprovalDecision{
		DecisionAt: time.Now().String(),
		Decision:   input.Decision,
		Comment:    input.Comment,
		TaskID:     input.TaskID,
	}
	if err := config.DB.Create(&decision).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, decision)
}

// PATCH /approval-decisions/:id
func UpdateApprovalDecision(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid approval decision id"})
		return
	}
	var decision entity.ApprovalDecision
	if err := config.DB.First(&decision, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Approval decision not found"})
		return
	}
	var input struct {
		Decision *string `json:"decision"`
		Comment  *string `json:"comment"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := make(map[string]interface{})
	if input.Decision != nil {
		updates["decision"] = *input.Decision
	}
	if input.Comment != nil {
		updates["comment"] = *input.Comment
	}
	if len(updates) > 0 {
		if err := config.DB.Model(&decision).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, decision)
}

// DELETE /approval-decisions/:id
func DeleteApprovalDecision(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid approval decision id"})
		return
	}
	if err := config.DB.Delete(&entity.ApprovalDecision{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Approval decision deleted successfully"})
}

// --- ApprovalRequirement Functions ---

// GET /approval-requirements
func GetApprovalRequirements(ctx *gin.Context) {
	var requirements []entity.ApprovalRequirement
	if err := config.DB.Preload("Scholarship").Preload("ApplicationDocuments").Find(&requirements).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, requirements)
}

// GET /approval-requirements/:id
func GetApprovalRequirementByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid approval requirement id"})
		return
	}
	var requirement entity.ApprovalRequirement
	if err := config.DB.Preload("Scholarship").Preload("ApplicationDocuments").First(&requirement, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Approval requirement not found"})
		return
	}
	ctx.JSON(http.StatusOK, requirement)
}

// POST /approval-requirements
func CreateApprovalRequirement(ctx *gin.Context) {
	var input struct {
		Name          string `json:"name" binding:"required"`
		Description   string `json:"description"`
		ScholarshipID uint   `json:"scholarship_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requirement := entity.ApprovalRequirement{
		Name:          input.Name,
		Description:   input.Description,
		ScholarshipID: input.ScholarshipID,
	}
	if err := config.DB.Create(&requirement).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, requirement)
}

// PATCH /approval-requirements/:id
func UpdateApprovalRequirement(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid approval requirement id"})
		return
	}
	var requirement entity.ApprovalRequirement
	if err := config.DB.First(&requirement, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Approval requirement not found"})
		return
	}
	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if len(updates) > 0 {
		if err := config.DB.Model(&requirement).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, requirement)
}

// DELETE /approval-requirements/:id
func DeleteApprovalRequirement(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid approval requirement id"})
		return
	}
	if err := config.DB.Delete(&entity.ApprovalRequirement{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Approval requirement deleted successfully"})
}
