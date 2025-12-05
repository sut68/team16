package entity

import "gorm.io/gorm"

type Screening struct {
	gorm.Model

	// ---------- User ----------
	UserID uint `json:"user_id" gorm:"not null"`
	User   User `json:"user" gorm:"foreignKey:UserID"`

	// ---------- Application ----------
	ApplicationID uint        `json:"application_id" gorm:"not null"`
	Application   Application `json:"application" gorm:"foreignKey:ApplicationID"`

	// ---------- Status Screening ----------
	StatusScreeningID uint            `json:"status_screening_id" gorm:"not null"`
	StatusScreening   StatusScreening `json:"status_screening" gorm:"foreignKey:StatusScreeningID"`

	// ---------- Criterial ----------
	/*CriterialID uint     `json:"criterial_id" gorm:"not null"`
	Criterial   Criterial `json:"criterial" gorm:"foreignKey:CriterialID"`*/

	// ---------- Reject Reason ----------
	RejectionReason *string `json:"rejection_reason"`
	//เก็บเหตุผลที่ไม่ผ่านการคัดกรอง - *ถ้าไม่ผ่านให้เก็บเหตุผล ถ้าผ่านให้เป็นค่าว่าง
}
