package approval

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"backend/entity"
)

// SeedApprovalDecisions creates sample decisions for a few pending tasks.
func SeedApprovalDecisions(db *gorm.DB) error {
	// Check if data already exists
	var count int64
	db.Model(&entity.ApprovalDecision{}).Count(&count)
	if count > 0 {
		return nil // Data already seeded
	}

	// Find an admin to make the decisions
	var admin entity.AdminProfile
	if err := db.First(&admin).Error; err != nil {
		return fmt.Errorf("could not find an admin profile for seeding decisions: %w", err)
	}

	// --- Decision 1: Approve but not yet qualified ---
	var task1 entity.ApprovalTask
	// Find a pending task that is not part of ApplicationScholarship 2, to avoid conflict
	err := db.Joins("JOIN application_documents ON application_documents.id = approval_tasks.document_id").
		Where("approval_tasks.status = ?", "pending").
		Where("application_documents.application_scholarship_id != ?", 2).
		First(&task1).Error

	if err == nil {
		decision1 := entity.ApprovalDecision{
			DecisionAt: time.Now(),
			Decision:   "approve",
			Comment:    "The document looks good. Approved.",
			TaskID:     task1.ID,
			AdminID:    admin.ID,
		}
		if err := db.Create(&decision1).Error; err != nil {
			return err
		}
		// Update parent task status
		db.Model(&task1).Update("status", "approved")
	} else {
		fmt.Println("Could not find a pending task for Decision 1, skipping.")
	}

	// --- Decision 2: Reject a task NOT related to APP-2 ---
	var taskToReject entity.ApprovalTask
	// Find another pending task that does not belong to ApplicationScholarship ID 2
	err = db.Joins("JOIN application_documents ON application_documents.id = approval_tasks.document_id").
		Where("approval_tasks.status = ?", "pending").
		Where("application_documents.application_scholarship_id != ?", 2).
		Where("approval_tasks.id != ?", task1.ID). // Ensure it's a different task
		First(&taskToReject).Error

	if err == nil {
		decision2 := entity.ApprovalDecision{
			DecisionAt: time.Now(),
			Decision:   "reject",
			Comment:    "The document is not valid or expired.",
			TaskID:     taskToReject.ID,
			AdminID:    admin.ID,
		}
		if err := db.Create(&decision2).Error; err != nil {
			return err
		}
		db.Model(&taskToReject).Update("status", "rejected")

		var doc entity.ApplicationDocument
		db.First(&doc, taskToReject.DocumentID)
		db.Model(&entity.ApplicationScholarship{}).Where("id = ?", doc.ApplicationScholarshipID).Update("status", "rejected")
	} else {
		fmt.Println("Could not find a task to reject that wasn't for APP-2, skipping rejection seed.")
	}

	// --- Decision 3: Request Change ---
	var taskToChange entity.ApprovalTask
	// Find another pending task for the request-change status
	err = db.Where("status = ?", "pending").
		Not("id = ?", task1.ID).
		Not("id = ?", taskToReject.ID).
		First(&taskToChange).Error

	if err == nil {
		decision3 := entity.ApprovalDecision{
			DecisionAt: time.Now(),
			Decision:   "request-change",
			Comment:    "Please upload a clearer version of this document.",
			TaskID:     taskToChange.ID,
			AdminID:    admin.ID,
		}
		if err := db.Create(&decision3).Error; err != nil {
			return err
		}
		db.Model(&taskToChange).Update("status", "request-change")
	} else {
		fmt.Println("Could not find a task for Request Change, skipping.")
	}


	// --- Explicitly approve #APP-2 (ApplicationScholarship ID 2) ---
	var appSch2 entity.ApplicationScholarship
	if err := db.First(&appSch2, 2).Error; err == nil {
		// Only process if it's not already qualified
		if appSch2.Status != "qualified" {
			var taskToDecide entity.ApprovalTask
			err := db.Joins("JOIN application_documents ON application_documents.id = approval_tasks.document_id").
				Where("application_documents.application_scholarship_id = ?", appSch2.ID).
				Where("approval_tasks.status = ?", "pending").
				First(&taskToDecide).Error

			if err == nil {
				decision := entity.ApprovalDecision{
					DecisionAt: time.Now(),
					Decision:   "approve",
					Comment:    "Force-seeded: All documents approved for #APP-2.",
					TaskID:     taskToDecide.ID,
					AdminID:    admin.ID,
				}
				db.Create(&decision)
			}

			// Approve all tasks for this application
			db.Exec("UPDATE approval_tasks SET status = 'approved' WHERE document_id IN (SELECT id FROM application_documents WHERE application_scholarship_id = ?)", appSch2.ID)

			// Update the main application status to 'qualified'
			db.Model(&appSch2).Update("status", "qualified")
		}
	}


	// --- Add 1 more Approved and Qualified example (if possible) ---
	var anotherTask entity.ApprovalTask
	err = db.Where("status = ?", "pending").First(&anotherTask).Error

	if err == nil {
		var doc entity.ApplicationDocument
		db.First(&doc, anotherTask.DocumentID)

		if doc.ApplicationScholarshipID != 2 { // Ensure it's not APP-2 again
			decision := entity.ApprovalDecision{
				DecisionAt: time.Now(),
				Decision:   "approve",
				Comment:    "Document verified. Moving to interview stage.",
				TaskID:     anotherTask.ID,
				AdminID:    admin.ID,
			}
			db.Create(&decision)

			db.Exec("UPDATE approval_tasks SET status = 'approved' WHERE document_id IN (SELECT id FROM application_documents WHERE application_scholarship_id = ?)", doc.ApplicationScholarshipID)

			db.Model(&entity.ApplicationScholarship{}).Where("id = ?", doc.ApplicationScholarshipID).Update("status", "qualified")
		}
	}


	return nil
}