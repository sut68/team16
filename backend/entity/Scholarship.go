package entity

import "gorm.io/gorm"

type Scholarship struct {
	gorm.Model
	ScholarshipName string `gorm:"not null" json:"scholarship_name"`
	Description     string `gorm:"not null" json:"description"`
	OpenDate        string `gorm:"not null" json:"open_date"`
	CloseDate       string `gorm:"not null" json:"close_date"`

	StatusscholarshipID uint              `json:"statusscholarship_id"`
	Statusscholarship   Statusscholarship `gorm:"foreignKey:StatusscholarshipID" json:"statusscholarship"`

	TypescholarshipID uint            `json:"typescholarship_id"`
	Typescholarship   Typescholarship `gorm:"foreignKey:TypescholarshipID" json:"typescholarship"`

	ApprovalRequirements    []ApprovalRequirement    `gorm:"foreignKey:ScholarshipID" json:"approval_requirements"`
	ApplicationScholarships []ApplicationScholarship `gorm:"foreignKey:ScholarshipID" json:"application_scholarships"`
}
