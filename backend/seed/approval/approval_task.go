package approval

import (
	"backend/entity"
	"fmt"

	"gorm.io/gorm"
)

// SeedApprovalTasks creates an ApprovalTask for any ApplicationDocument that doesn't have one yet.
func SeedApprovalTasks(db *gorm.DB) error {
	var documents []entity.ApplicationDocument
	if err := db.Find(&documents).Error; err != nil {
		return fmt.Errorf("failed to find application documents: %w", err)
	}

	for _, doc := range documents {
		var task entity.ApprovalTask
		err := db.Where("document_id = ?", doc.ID).First(&task).Error

		// If a task for this document doesn't exist, create one.
		if err != nil && err == gorm.ErrRecordNotFound {
			newTask := entity.ApprovalTask{
				Status:     "pending", // Use lowercase as seen in the controller
				DocumentID: doc.ID,
			}
			if err := db.Create(&newTask).Error; err != nil {
				return fmt.Errorf("failed to create approval task for document ID %d: %w", doc.ID, err)
			}
		} else if err != nil {
			// Handle other potential errors
			return fmt.Errorf("failed to check for existing approval task for document ID %d: %w", doc.ID, err)
		}
		// If a task already exists, do nothing.
	}

	return nil
}