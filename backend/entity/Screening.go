package entity

import "gorm.io/gorm"

type Screening struct {
	gorm.Model

	// ---------- User ----------

	AdminProfileID uint         `json:"admin_profile_id" gorm:"not null" valid:"required"`
	AdminProfile   AdminProfile `json:"admin_profile" gorm:"foreignKey:AdminProfileID"`

	// ---------- Status Screening ----------
	StatusScreeningID uint            `json:"status_screening_id" gorm:"not null" valid:"required"`
	StatusScreening   StatusScreening `json:"status_screening" gorm:"foreignKey:StatusScreeningID"`

	// ---------- Reject Reason ----------
	RejectionReason *string `json:"rejection_reason" valid:"optional,stringlength(1|100)~You must provide a rejection reason with length between 1 and 100 characters"`
	
	// ---------- Application Scholarship ----------
	ApplicationScholarshipID uint                   `json:"application_scholarship_id" gorm:"not null" valid:"required"`
	ApplicationScholarship   ApplicationScholarship `json:"application_scholarship" gorm:"foreignKey:ApplicationScholarshipID"`
	
}
