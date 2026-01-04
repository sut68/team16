package entity

import "gorm.io/gorm"

type Screening struct {
	gorm.Model

	// ---------- User ----------

	AdminProfileID uint         `json:"admin_profile_id" gorm:"not null" valid:"required~AdminProfileID is required"`
	AdminProfile   AdminProfile `json:"admin_profile" gorm:"foreignKey:AdminProfileID" valid:"-"`

	// ---------- Status Screening ----------
	StatusScreeningID uint            `json:"status_screening_id" gorm:"not null" valid:"required~StatusScreeningID is required"`
	StatusScreening   StatusScreening `json:"status_screening" gorm:"foreignKey:StatusScreeningID" valid:"-"`

	// ---------- Reject Reason ----------
	RejectionReason *string `json:"rejection_reason" valid:"maxstringlength(100)~Rejection Reason must be less than 100 characters"`

	ApplicationScholarshipID uint                   `json:"application_scholarship_id" gorm:"not null" valid:"required~ApplicationScholarshipID is required"`
	ApplicationScholarship   ApplicationScholarship `json:"application_scholarship" gorm:"foreignKey:ApplicationScholarshipID" valid:"-"`
}
