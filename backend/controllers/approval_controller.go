package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/entity"
	"backend/storage"
	"backend/validators"
)

func GetApprovalTasks(ctx *gin.Context) {

	var tasks []entity.ApprovalTask

	if err := config.DB.
		Preload("ApplicationDocument.ApplicationScholarship.Application.StudentProfile.Major").
		Preload("ApplicationDocument.ApplicationScholarship.Application.Semaster").
		Preload("ApplicationDocument.ApplicationScholarship.Scholarship.ApprovalRequirements.Requirement").
		Preload("ApplicationDocument.ApplicationScholarship.ApplicationDocuments").
		Preload("ApprovalDecisions.Admin").
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
		Preload("ApplicationDocument.ApplicationScholarship.Application.StudentProfile.Major").
		Preload("ApplicationDocument.ApplicationScholarship.Application.Semaster").
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
		Status *string `json:"status" valid:"optional,in(pending|approved|rejected|request-change)~Invalid status"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateStruct(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"error":   err.Error(),
		})
		return
	}

	tx := config.DB.Begin()

	updates := make(map[string]interface{})
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if len(updates) > 0 {
		if err := tx.Model(&task).Updates(updates).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "commit transaction failed"})
		return
	}

	if err := config.DB.Preload("ApplicationDocument").Preload("ApprovalDecisions").First(&task, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "reload failed"})
		return
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
		Preload("Admin").
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
	// --- 1. Get Input & Handle File Upload ---
	appScholarshipIDStr := ctx.PostForm("application_scholarship_id")
	appScholarshipID, _ := strconv.ParseUint(appScholarshipIDStr, 10, 64)

	uploadedBy := ctx.PostForm("uploaded_by") // Assuming student profile ID is sent

	file, err := ctx.FormFile("document")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	var filePath string
	var uniqueFileName string

	// Try to upload to MinIO first, fallback to local storage
	if storage.IsConfigured() {
		// Upload to MinIO
		objectKey, publicURL, err := storage.Client.UploadFile(file, "application")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to storage: " + err.Error()})
			return
		}
		uniqueFileName = objectKey
		filePath = publicURL
	} else {
		// Fallback to local storage
		uniqueFileName = fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
		filePath = fmt.Sprintf("uploads/application/%s", uniqueFileName)

		// Ensure upload directory exists
		if _, err := os.Stat("uploads/application"); os.IsNotExist(err) {
			os.MkdirAll("uploads/application", 0755)
		}

		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
			return
		}
	}

	// --- 2. Create ApplicationDocument record ---
	document := entity.ApplicationDocument{
		FileName:                 uniqueFileName,
		FilePath:                 filePath,
		FileType:                 file.Header.Get("Content-Type"),
		UploadedBy:               uploadedBy,
		ApplicationScholarshipID: uint(appScholarshipID),
	}

	if err := validators.ValidateStruct(&document); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"error":   err.Error(),
		})
		return
	}

	tx := config.DB.Begin()

	if err := tx.Create(&document).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document record: " + err.Error()})
		return
	}

	// --- 3. Create associated ApprovalTask ---
	var existingTask entity.ApprovalTask

	errTask := tx.Joins("JOIN application_documents ON application_documents.id = approval_tasks.document_id").
		Where("application_documents.application_scholarship_id = ?", appScholarshipID).
		First(&existingTask).Error

	if errTask == nil {
		if err := tx.Model(&existingTask).Updates(map[string]interface{}{
			"document_id": document.ID,
			"status":      "pending",
		}).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update existing task: " + err.Error()})
			return
		}

		// Reset ApplicationScholarship status to pending when new document is uploaded
		// This ensures the status chain is properly reset for re-review
		if err := tx.Model(&entity.ApplicationScholarship{}).
			Where("id = ?", appScholarshipID).
			Update("status", "pending").Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application status: " + err.Error()})
			return
		}

	} else {

		newTask := entity.ApprovalTask{
			Status:     "pending",
			DocumentID: document.ID,
		}

		if err := tx.Create(&newTask).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create approval task: " + err.Error()})
			return
		}
	}

	// --- 4. Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, document)
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
	if err := config.DB.Preload("ApprovalTask").Preload("Admin").Find(&decisions).Error; err != nil {
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
	if err := config.DB.Preload("ApprovalTask").Preload("Admin").First(&decision, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Approval decision not found"})
		return
	}
	ctx.JSON(http.StatusOK, decision)
}

// POST /approval-decisions
func CreateApprovalDecision(ctx *gin.Context) {
	var input struct {
		Decision string `json:"decision"`
		Comment  string `json:"comment"`
		TaskID   uint   `json:"task_id"`
		AdminID  uint   `json:"admin_id"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	decision := entity.ApprovalDecision{
		DecisionAt: time.Now(),
		Decision:   input.Decision,
		Comment:    input.Comment,
		TaskID:     input.TaskID,
		AdminID:    input.AdminID,
	}

	if err := validators.ValidateStruct(&decision); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"error":   err.Error(),
		})
		return
	}

	tx := config.DB.Begin()

	if err := tx.Create(&decision).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update ApprovalTask status based on decision
	taskStatus := "pending"
	switch input.Decision {
	case "approve":
		taskStatus = "approved"
	case "reject":
		taskStatus = "rejected"
	case "request-change":
		taskStatus = "request-change"
	}

	if err := tx.Model(&entity.ApprovalTask{}).Where("id = ?", input.TaskID).Update("status", taskStatus).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
		return
	}

	// If decision is "approve", check if all documents for this application are approved
	// If so, update ApplicationScholarship.status to "qualified"
	if input.Decision == "approve" {
		var task entity.ApprovalTask
		if err := tx.Preload("ApplicationDocument").First(&task, input.TaskID).Error; err == nil {
			appScholarshipID := task.ApplicationDocument.ApplicationScholarshipID

			// Count documents that are NOT approved for this application
			var pendingCount int64
			tx.Model(&entity.ApprovalTask{}).
				Joins("JOIN application_documents ON application_documents.id = approval_tasks.document_id").
				Where("application_documents.application_scholarship_id = ?", appScholarshipID).
				Where("approval_tasks.status != ?", "approved").
				Count(&pendingCount)

			// If all documents are approved (pendingCount == 0), set status to "qualified"
			if pendingCount == 0 {
				if err := tx.Model(&entity.ApplicationScholarship{}).
					Where("id = ?", appScholarshipID).
					Update("status", "qualified").Error; err != nil {
					tx.Rollback()
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application status"})
					return
				}
			}
		}
	}

	// If decision is "reject", update ApplicationScholarship.status to "rejected"
	if input.Decision == "reject" {
		var task entity.ApprovalTask
		if err := tx.Preload("ApplicationDocument").First(&task, input.TaskID).Error; err == nil {
			appScholarshipID := task.ApplicationDocument.ApplicationScholarshipID
			tx.Model(&entity.ApplicationScholarship{}).
				Where("id = ?", appScholarshipID).
				Update("status", "rejected")
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	ctx.JSON(http.StatusCreated, decision)
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
	var requirement entity.ApprovalRequirement
	if err := ctx.ShouldBindJSON(&requirement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateStruct(&requirement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"error":   err.Error(),
		})
		return
	}

	if err := config.DB.Create(&requirement).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, requirement)
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
