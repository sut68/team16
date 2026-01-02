package interview

import (
	"backend/entity"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func SeedInterviewRounds(db *gorm.DB) error {
	// Check if data already exists
	var count int64
	db.Model(&entity.InterviewRound{}).Count(&count)
	if count > 0 {
		return nil // Data already seeded
	}

	// Find all scholarships that are "qualified" to have interviews
	var qualifiedApps []entity.ApplicationScholarship
	if err := db.Where("status = ?", "qualified").Preload("Scholarship").Find(&qualifiedApps).Error; err != nil {
		return fmt.Errorf("error finding qualified applications: %w", err)
	}

	if len(qualifiedApps) == 0 {
		fmt.Println("No qualified applications found for seeding interview rounds.")
		return nil
	}

	// Find common entities needed for the rounds
	var adminProfile entity.AdminProfile
	if err := db.First(&adminProfile).Error; err != nil {
		return fmt.Errorf("no admin profile found: %w", err)
	}

	var location entity.Location
	if err := db.First(&location).Error; err != nil {
		return fmt.Errorf("no location found: %w", err)
	}

	var interviewMode entity.InterviewMode
	if err := db.First(&interviewMode).Error; err != nil {
		return fmt.Errorf("no interview mode found: %w", err)
	}

	// Find some interviewers to assign to slots
	var interviewers []entity.Interviewer
	if err := db.Limit(2).Find(&interviewers).Error; err != nil || len(interviewers) == 0 {
		return fmt.Errorf("could not find interviewers to seed slots: %w", err)
	}

	// --- Create a round for each qualified application ---
	for i, app := range qualifiedApps {
		round := entity.InterviewRound{
			Name:            fmt.Sprintf("รอบสัมภาษณ์ทุน %s", app.Scholarship.ScholarshipName),
			Description:     fmt.Sprintf("รอบสัมภาษณ์สำหรับผู้ผ่านการคัดเลือกเอกสาร %s", app.Scholarship.ScholarshipName),
			StartDateTime:   time.Now().AddDate(0, 0, 7+(i*2)), // Stagger start dates
			EndDateTime:     time.Now().AddDate(0, 0, 7+(i*2)).Add(time.Hour * 3), // 3 hours of interviews
			SlotDuration:    30, // 30 minutes
			ScholarshipID:   app.ScholarshipID,
			AdminProfileID:  adminProfile.ID,
			LocationID:      &location.ID,
			InterviewModeID: &interviewMode.ID,
			MeetingLink:     "https://meet.google.com/seed-data-link",
		}

		// Use a transaction to create round, slots, and assignments
		err := db.Transaction(func(tx *gorm.DB) error {
			// 1. Create the InterviewRound
			if err := tx.Create(&round).Error; err != nil {
				return fmt.Errorf("failed to create round for scholarship %d: %w", app.ScholarshipID, err)
			}

			// 2. Generate and create Slots
			var slots []entity.Slot
			duration := time.Duration(round.SlotDuration) * time.Minute
			for t := round.StartDateTime; t.Before(round.EndDateTime); t = t.Add(duration) {
				slots = append(slots, entity.Slot{
					InterviewRoundID: round.ID,
					StartTime:        t,
					EndTime:          t.Add(duration),
					Capacity:         1,
					BookCount:        0,
					Status:           "Available",
					IsBooked:         false,
				})
			}

			if len(slots) > 0 {
				if err := tx.Create(&slots).Error; err != nil {
					return fmt.Errorf("failed to create slots for round %d: %w", round.ID, err)
				}
			}

			// 3. Create InterviewerSlot assignments
			var interviewerSlots []entity.InterviewerSlot
			for _, slot := range slots {
				// Assign both interviewers to each slot
				for _, interviewer := range interviewers {
					interviewerSlots = append(interviewerSlots, entity.InterviewerSlot{
						SlotID:        slot.ID,
						InterviewerID: interviewer.ID,
					})
				}
			}

			if len(interviewerSlots) > 0 {
				if err := tx.Create(&interviewerSlots).Error; err != nil {
					return fmt.Errorf("failed to assign interviewers for round %d: %w", round.ID, err)
				}
			}

			return nil
		})

		if err != nil {
			return err // If transaction fails, return the error
		}
	}

	return nil
}