package entity

import "gorm.io/gorm"

type Screening struct {
	gorm.Model

	// ---------- User ----------

	AdminProfileID uint         `json:"admin_profile_id" gorm:"not null"`
	AdminProfile   AdminProfile `json:"admin_profile" gorm:"foreignKey:AdminProfileID"`

	// ---------- Application ----------
	ApplicationID uint        `json:"application_id" gorm:"not null"`
	Application   Application `json:"application" gorm:"foreignKey:ApplicationID"`

	// ---------- Status Screening ----------
	StatusScreeningID uint            `json:"status_screening_id" gorm:"not null"`
	StatusScreening   StatusScreening `json:"status_screening" gorm:"foreignKey:StatusScreeningID"`

	// ---------- Reject Reason ----------
	RejectionReason *string `json:"rejection_reason"`

	ApplicationScholarshipID uint                   `json:"application_scholarship_id" gorm:"not null"`
	ApplicationScholarship   ApplicationScholarship `json:"application_scholarship" gorm:"foreignKey:ApplicationScholarshipID"`
}