package interview

import (
	"log"
	"backend/entity"

	"gorm.io/gorm"
)

func SeedInterviewModes(db *gorm.DB) {
	modes := []entity.InterviewMode{
		{Name: "Onsite"},
		{Name: "Online"},
	}

	for _, mode := range modes {
		var existingMode entity.InterviewMode
		// Check if the mode already exists
		if err := db.Where("name = ?", mode.Name).First(&existingMode).Error; err == gorm.ErrRecordNotFound {
			// Create the mode if it does not exist
			if err := db.Create(&mode).Error; err != nil {
				// Handle error, e.g., log it
				log.Println("Failed to seed interview mode: " + err.Error())
			}
		}
	}
}
