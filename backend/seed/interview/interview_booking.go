package interview

import (
	"backend/entity"
	"fmt"

	"gorm.io/gorm"
)

// SeedInterviewBookings creates a sample booking for an available slot.
func SeedInterviewBookings(db *gorm.DB) error {
	// Check if data already exists
	var count int64
	db.Model(&entity.IntervieweBooking{}).Count(&count)
	if count > 0 {
		return nil // Data already seeded
	}

	// 1. Find a student application that is 'qualified' for an interview
	var appScholarship entity.ApplicationScholarship
	if err := db.Where("status = ?", "qualified").
		Preload("Application.StudentProfile").
		First(&appScholarship).Error; err != nil {
		return fmt.Errorf("no qualified application found to seed a booking: %w", err)
	}

	// 2. Find an available interview slot for that scholarship
	var availableSlot entity.Slot
	if err := db.Joins("JOIN interview_rounds ON interview_rounds.id = slots.interview_round_id").
		Where("interview_rounds.scholarship_id = ? AND slots.status = ?", appScholarship.ScholarshipID, "Available").
		First(&availableSlot).Error; err != nil {
		return fmt.Errorf("no available slot found for scholarship ID %d to seed a booking: %w", appScholarship.ScholarshipID, err)
	}

	// 3. Create the booking and update the slot in a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// Create the booking
		booking := entity.IntervieweBooking{
			Status:                   "booked",
			SlotID:                   availableSlot.ID,
			ApplicationScholarshipID: appScholarship.ID,
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		// Update the slot's status and book count
		if err := tx.Model(&availableSlot).Updates(map[string]interface{}{
			"book_count": availableSlot.BookCount + 1,
			"status":     "Booked",
			"is_booked":  true,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
