package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/entity"
)

// GetAllInterviewRounds godoc
func GetAllInterviewRounds(c *gin.Context) {
	var rounds []entity.InterviewRound
	if err := config.DB.
		Preload("Scholarship").
		Preload("AdminProfile").
		Preload("Slots").
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
		Preload("Slots").
		Preload("Slots.InterviewerSlots.Interviewer").
		Preload("Slots.InterviewBooking").
		First(&round, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview round not found"})
		return
	}
	c.JSON(http.StatusOK, round)
}

// CreateInterviewRound godoc
func CreateInterviewRound(c *gin.Context) {
	var round entity.InterviewRound
	if err := c.ShouldBindJSON(&round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&round).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, round)
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

	// TODO: Add validation to ensure the slot is not already booked

	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
