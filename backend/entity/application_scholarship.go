package entity

import "gorm.io/gorm"

type ApplicationScholarship struct {
	gorm.Model
	Status string `json:"status"`

	ApplicationID uint        `json:"application_id"`
	Application   Application `gorm:"foreignKey:ApplicationID" json:"application"`

	ScholarshipID uint        `json:"scholarship_id"`
	Scholarship   Scholarship `gorm:"foreignKey:ScholarshipID" json:"scholarship"`

	ApplicationDocuments []ApplicationDocument `gorm:"foreignKey:ApplicationScholarshipID" json:"application_documents"`
	IntervieweBookings  []IntervieweBooking   `gorm:"foreignKey:ApplicationScholarshipID" json:"interviewe_bookings"`
}
