package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/entity"
)

func GetApprovalTasks(ctx *gin.Context) {

	var tasks []entity.ApprovalTask

	if err := config.DB.
		Preload("Admin").
		Preload("ApplicationDocument.ApplicationScholarship.Application.StudentProfile.Major").
		Preload("ApplicationDocument.ApplicationScholarship.Scholarship.ApprovalRequirements.Requirement").
		Preload("ApplicationDocument.ApplicationScholarship.ApplicationDocuments").
		Preload("ApprovalDecisions").
		Find(&tasks).Error; err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return

	}

	// DEBUGGING: Print the fetched tasks to console

	jsonData, err := json.MarshalIndent(tasks, "", "  ")

	if err != nil {

		fmt.Println("Error marshalling tasks for debugging:", err)

	} else {

		fmt.Println("--- DEBUG: Fetched Approval Tasks ---")

		fmt.Println(string(jsonData))

		fmt.Println("--- END DEBUG ---")

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
		Preload("ApplicationDocument.ApplicationScholarship.Application.StudentProfile.Major").
		Preload("ApplicationDocument.ApplicationScholarship.Scholarship.ApprovalRequirements.Requirement").
		Preload("ApprovalDecisions").
		First(&task, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Approval task not found"})
		return
	}
	ctx.JSON(http.StatusOK, task)
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
		Joins("JOIN application_scholarships ON application_scholarships.id = application_documents.application_scholarship_id").
		Joins("JOIN applications ON applications.id = application_scholarships.application_id").
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

// GET /application-documents
func GetApplicationDocuments(ctx *gin.Context) {
	var documents []entity.ApplicationDocument
	if err := config.DB.Preload("ApplicationScholarship.Application.StudentProfile").Find(&documents).Error; err != nil {
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
	if err := config.DB.Preload("ApplicationScholarship.Application.StudentProfile").First(&document, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Application document not found"})
		return
	}
	ctx.JSON(http.StatusOK, document)
}

// POST /application-documents
func CreateApplicationDocument(ctx *gin.Context) {
	// Start a new database transaction
	tx := config.DB.Begin()

	// --- 1. Get Input & Handle File Upload ---
	appScholarshipIDStr := ctx.PostForm("application_scholarship_id")
	if appScholarshipIDStr == "" {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "application_scholarship_id is a required field"})
		return
	}
	appScholarshipID, _ := strconv.ParseUint(appScholarshipIDStr, 10, 64)

	uploadedBy := ctx.PostForm("uploaded_by") // Assuming student profile ID is sent
	if uploadedBy == "" {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "uploaded_by is a required field"})
		return
	}

	file, err := ctx.FormFile("document")
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Generate a unique filename and path
	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	filePath := fmt.Sprintf("uploads/%s", uniqueFileName)

	// Save the file
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	// --- 2. Create ApplicationDocument record ---
	document := entity.ApplicationDocument{
		FileName:                 uniqueFileName,
		FilePath:                 filePath,
		FileType:                 file.Header.Get("Content-Type"),
		UploadedBy:               uploadedBy,
		ApplicationScholarshipID: uint(appScholarshipID),
	}

	if err := tx.Create(&document).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document record: " + err.Error()})
		return
	}

	// --- 3. Create associated ApprovalTask ---
	task := entity.ApprovalTask{
		Status:     "pending", // Initial status for admin review
		AdminID:    1,         // Assumption: A default admin (ID 1) for initial assignment
		DocumentID: document.ID,
	}

	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create approval task: " + err.Error()})
		return
	}

	// Update the status of the ApplicationScholarship to "pending"
	var appScholarship entity.ApplicationScholarship
	if err := tx.First(&appScholarship, document.ApplicationScholarshipID).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find parent application scholarship: " + err.Error()})
		return
	}
	if err := tx.Model(&appScholarship).Update("status", "pending").Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application scholarship status: " + err.Error()})
		return
	}

	// --- 4. Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed: " + err.Error()})
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

// GET /approval-requirements
func GetApprovalRequirements(ctx *gin.Context) {
	var requirements []entity.ApprovalRequirement
	if err := config.DB.Preload("Scholarship").Preload("Requirement").Find(&requirements).Error; err != nil {
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
	if err := config.DB.Preload("Scholarship").Preload("Requirement").First(&requirement, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Approval requirement not found"})
		return
	}
	ctx.JSON(http.StatusOK, requirement)
}

// POST /approval-requirements
func CreateApprovalRequirement(ctx *gin.Context) {
	var input struct {
		ScholarshipID uint `json:"scholarship_id" binding:"required"`
		RequirementID uint `json:"requirement_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requirement := entity.ApprovalRequirement{
		ScholarshipID: input.ScholarshipID,
		RequirementID: input.RequirementID,
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
