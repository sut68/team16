package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/entity"
)

// GetAllInterviewRounds godoc
func GetAllInterviewRounds(c *gin.Context) {
	var rounds []entity.InterviewRound
	if err := config.DB.
		Preload("Scholarship.Statusscholarship").
		Preload("Scholarship.Typescholarship").
		Preload("AdminProfile").
		Preload("InterviewMode").
		Preload("Location").
		Preload("Slots.InterviewerSlots.Interviewer").
		Find(&rounds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rounds)
}

// GetInterviewRoundByID godoc
func GetInterviewRoundByID(c *gin.Context) {
	id := c.Param("id")
	var round entity.InterviewRound
	if err := config.DB.
		Preload("Scholarship").
		Preload("AdminProfile").
		Preload("InterviewMode").
		Preload("Location").
		Preload("Slots.InterviewerSlots.Interviewer").
		Preload("Slots.IntervieweBookings.ApplicationScholarship.Application.StudentProfile.Major").
		First(&round, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview round not found"})
		return
	}
	c.JSON(http.StatusOK, round)
}

// CreateInterviewRound godoc
// CreateInterviewRound godoc
func CreateInterviewRound(c *gin.Context) {
	var input struct {
		entity.InterviewRound
		InterviewerIDs []uint `json:"interviewer_ids"`
		Slots          []struct {
			Status string `json:"status"`
		} `json:"slots"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	// 1. Create the InterviewRound
	round := input.InterviewRound
	if err := tx.Create(&round).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interview round: " + err.Error()})
		return
	}

	// 2. Calculate and create slots using status from payload
	startTime := round.StartDateTime
	endTime := round.EndDateTime
	duration := time.Duration(round.SlotDuration) * time.Minute
	var slots []entity.Slot

	i := 0
	for currentTime := startTime; currentTime.Add(duration).Before(endTime) || currentTime.Add(duration).Equal(endTime); currentTime = currentTime.Add(duration) {
		slotStatus := "Available" // Default
		if i < len(input.Slots) {
			slotStatus = input.Slots[i].Status
		}

		slot := entity.Slot{
			InterviewRoundID: round.ID,
			StartTime:        currentTime,
			EndTime:          currentTime.Add(duration),
			Capacity:         1,
			BookCount:        0,
			Status:           slotStatus,
		}
		slots = append(slots, slot)
		i++
	}

	if len(slots) > 0 {
		if err := tx.Create(&slots).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create slots: " + err.Error()})
			return
		}
	}

	// 3. Assign interviewers to each slot
	if len(input.InterviewerIDs) > 0 && len(slots) > 0 {
		var interviewerSlots []entity.InterviewerSlot
		for _, slot := range slots {
			for _, interviewerID := range input.InterviewerIDs {
				interviewerSlot := entity.InterviewerSlot{
					SlotID:        slot.ID,
					InterviewerID: interviewerID,
				}
				interviewerSlots = append(interviewerSlots, interviewerSlot)
			}
		}
		if len(interviewerSlots) > 0 {
			if err := tx.Create(&interviewerSlots).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign interviewers: " + err.Error()})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed: " + err.Error()})
		return
	}

	// Refetch the created round
	var createdRound entity.InterviewRound
	if err := config.DB.
		Preload("Scholarship").
		Preload("AdminProfile").
		Preload("Slots.InterviewerSlots.Interviewer").
		First(&createdRound, round.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created round: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdRound)
}

// UpdateInterviewRound godoc
func UpdateInterviewRound(c *gin.Context) {
	id := c.Param("id")
	var round entity.InterviewRound
	if err := config.DB.First(&round, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview round not found"})
		return
	}

	var input struct {
		entity.InterviewRound
		Slots []struct {
			ID     uint   `json:"id"`
			Status string `json:"status"`
		} `json:"slots"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	// 1. Update the round details
	if err := tx.Model(&round).Updates(input.InterviewRound).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update interview round: " + err.Error()})
		return
	}

	// 2. Update individual slot statuses if provided
	if len(input.Slots) > 0 {
		for _, slotUpdate := range input.Slots {
			var slotToUpdate entity.Slot
			if err := tx.First(&slotToUpdate, slotUpdate.ID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"error": "Slot with ID " + c.Param("id") + " not found"})
				return
			}

			// Do not allow changing status of a booked slot
			if slotToUpdate.IsBooked {
				continue
			}

			if err := tx.Model(&slotToUpdate).Update("status", slotUpdate.Status).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update slot status: " + err.Error()})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed: " + err.Error()})
		return
	}

	// Refetch the round to return the updated data
	var updatedRound entity.InterviewRound
	if err := config.DB.
		Preload("Scholarship").
		Preload("AdminProfile").
		Preload("Slots.InterviewerSlots.Interviewer").
		First(&updatedRound, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated round: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedRound)
}

// DeleteInterviewRound godoc
func DeleteInterviewRound(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&entity.InterviewRound{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Interview round deleted"})
}

// --- Interviewer Handlers ---

// GetAllInterviewers godoc
func GetAllInterviewers(c *gin.Context) {
	var interviewers []entity.Interviewer
	if err := config.DB.Find(&interviewers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, interviewers)
}

// CreateInterviewer godoc
func CreateInterviewer(c *gin.Context) {
	var input entity.Interviewer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use FirstOrCreate to find an interviewer by email or create a new one.
	// The 'where' clause specifies the unique field to look for.
	// The struct passed as the second argument contains the data for creation if not found.
	if err := config.DB.Where(entity.Interviewer{Email: input.Email}).
		Attrs(entity.Interviewer{
			InterviewerFirstname: input.InterviewerFirstname,
			InterviewerLastname:  input.InterviewerLastname,
		}).
		FirstOrCreate(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find or create interviewer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, input) // Return the found or newly created interviewer
}

// --- InterviewBooking Handlers ---

// CreateInterviewBooking godoc
func CreateInterviewBooking(c *gin.Context) {
	var booking entity.IntervieweBooking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	// 1. Check if the slot is available
	var slot entity.Slot
	if err := tx.First(&slot, booking.SlotID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Slot not found"})
		return
	}

	if slot.Status != "Available" {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Slot is not available for booking"})
		return
	}

	if slot.BookCount >= slot.Capacity {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Slot is fully booked"})
		return
	}

	// 2. Create the booking
	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// 3. Increment the slot's book_count and set is_booked to true
	if err := tx.Model(&slot).Updates(map[string]interface{}{"book_count": slot.BookCount + 1, "is_booked": true}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update slot status"})
		return
	}

	// 4. Update ApplicationScholarship status to 'interview_scheduled'
	if err := tx.Model(&entity.ApplicationScholarship{}).
		Where("id = ?", booking.ApplicationScholarshipID).
		Update("status", "interview_scheduled").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application status"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// GetStudentBookings godoc
func GetStudentBookings(c *gin.Context) {
	studentProfileID := c.Param("student_profile_id")
	var bookings []entity.IntervieweBooking

	// This query is a bit complex. It joins through ApplicationScholarship and Application
	// to find bookings related to a student profile.
	if err := config.DB.
		Joins("JOIN application_scholarships ON application_scholarships.id = interviewe_bookings.application_scholarship_id").
		Joins("JOIN applications ON applications.id = application_scholarships.application_id").
		Where("applications.student_profile_id = ?", studentProfileID).
		Preload("Slot").
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// DeleteInterviewBooking godoc
func DeleteInterviewBooking(c *gin.Context) {
	id := c.Param("id")
	var booking entity.IntervieweBooking

	// Find the booking to be deleted
	if err := config.DB.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	// 1. Delete the booking itself
	if err := tx.Delete(&booking).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	// 2. Find the associated slot
	var slot entity.Slot
	if err := tx.First(&slot, booking.SlotID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Associated slot not found"})
		return
	}

	// 3. Decrement the book_count and update the status
	newBookCount := uint(0)
	if slot.BookCount > 0 {
		newBookCount = slot.BookCount - 1
	}

	isBooked := newBookCount > 0

	if err := tx.Model(&slot).Updates(map[string]interface{}{
		"book_count": newBookCount,
		"is_booked":  isBooked,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update slot status"})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interview booking cancelled and slot is now available"})
}

// GetAllLocations godoc
func GetAllLocations(c *gin.Context) {
	var locations []entity.Location
	if err := config.DB.Find(&locations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locations)
}

// GetAllInterviewModes godoc
func GetAllInterviewModes(c *gin.Context) {
	var modes []entity.InterviewMode
	if err := config.DB.Find(&modes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, modes)
}

