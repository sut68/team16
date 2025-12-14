package entity

import "gorm.io/gorm"

type IntervieweBooking struct {
	gorm.Model
	BookingDate string `gorm:"not null" json:"booking_date"`
	BookingTime string `gorm:"not null" json:"booking_time"`
	Status 	string `gorm:"not null" json:"status"`

	SlotID uint   `json:"slot_id"`
	Slot   Slot   `gorm:"foreignKey:SlotID" json:"slot"`

	ApplicationScholarshipID uint                   `json:"application_scholarship_id"`
	ApplicationScholarship   ApplicationScholarship `gorm:"foreignKey:ApplicationScholarshipID" json:"application_scholarship"`


}