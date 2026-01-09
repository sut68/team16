package entity

import "gorm.io/gorm"

type IntervieweBooking struct {
	gorm.Model
	Status string `gorm:"not null" json:"status"`

	// Who booked this slot?
	BookedByRole    string       `gorm:"type:varchar(20);not null;default:'student'" json:"booked_by_role"` // 'student' or 'admin'
	BookedByAdminID *uint        `json:"booked_by_admin_id,omitempty"`
	AdminProfile    AdminProfile `gorm:"foreignKey:BookedByAdminID" json:"admin_profile,omitempty"`

	SlotID uint `json:"slot_id"`
	Slot   Slot `gorm:"foreignKey:SlotID" json:"slot"`

	ApplicationScholarshipID uint                   `json:"application_scholarship_id"`
	ApplicationScholarship   ApplicationScholarship `gorm:"foreignKey:ApplicationScholarshipID" json:"application_scholarship"`
}