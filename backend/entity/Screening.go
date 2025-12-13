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
	//เก็บเหตุผลที่ไม่ผ่านการคัดกรอง - *ถ้าไม่ผ่านให้เก็บเหตุผล ถ้าผ่านให้เป็นค่าว่าง

	ScholarshipID uint        `json:"scholarship_id" gorm:"not null"`
	Scholarship   Scholarship `json:"scholarship" gorm:"foreignKey:ScholarshipID"`
}
