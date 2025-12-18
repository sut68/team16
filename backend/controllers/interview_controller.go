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
		Preload("Slots").
		Preload("Slots.InterviewerSlots.Interviewer").
		Preload("Slots.IntervieweBookings").
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

	// 2. Calculate and create slots
	startTime := round.StartDateTime
	endTime := round.EndDateTime

	duration := time.Duration(round.SlotDuration) * time.Minute
	var slots []entity.Slot

	for currentTime := startTime; currentTime.Add(duration).Before(endTime) || currentTime.Add(duration).Equal(endTime); currentTime = currentTime.Add(duration) {
		
		slot := entity.Slot{
			InterviewRoundID: round.ID,
			StartTime:        currentTime,
			EndTime:          currentTime.Add(duration),
			Capacity:         1,
			BookCount:        0,
			Status:           "Available",
		}
		slots = append(slots, slot)
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
		if err := tx.Create(&interviewerSlots).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign interviewers: " + err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed: " + err.Error()})
		return
	}

	// Refetch the created round ... (ส่วนนี้เหมือนเดิม)
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

	var input entity.InterviewRound
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&round).Updates(input)
	c.JSON(http.StatusOK, round)
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
	var interviewer entity.Interviewer
	if err := c.ShouldBindJSON(&interviewer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&interviewer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, interviewer)
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
	
	// Find the booking
	if err := config.DB.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	// Delete the booking
	if err := tx.Delete(&booking).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	// Update the slot's is_booked status
	if err := tx.Model(&entity.Slot{}).Where("id = ?", booking.SlotID).Update("is_booked", false).Error; err != nil {
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

