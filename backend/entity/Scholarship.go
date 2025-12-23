package entity

import "gorm.io/gorm"

type Scholarship struct {
	gorm.Model
	ScholarshipName string `gorm:"not null" json:"scholarship_name" valid:"required~Scholarship name is required,stringlength(2|255)~Scholarship name must be 2-255 characters"`
	Description     string `gorm:"not null" json:"description" valid:"required~Description is required,stringlength(5|1000)~Description must be 5-1000 characters"`
	OpenDate        string `gorm:"not null" json:"open_date" valid:"required~Open date is required"`
	CloseDate       string `gorm:"not null" json:"close_date" valid:"required~Close date is required"`

	StatusscholarshipID uint              `json:"statusscholarship_id" valid:"required~Status scholarship is required"`
	Statusscholarship   Statusscholarship `gorm:"foreignKey:StatusscholarshipID" json:"statusscholarship"`
 
	TypescholarshipID uint            `json:"typescholarship_id" valid:"required~Type scholarship is required"`
	Typescholarship   Typescholarship `gorm:"foreignKey:TypescholarshipID" json:"typescholarship" valid:"-"`

	SemasterID uint     `json:"semaster_id" valid:"required~Semaster is required"`
	Semaster   Semaster `gorm:"foreignKey:SemasterID" json:"semaster" valid:"-"`

	SponsorID uint `json:"sponsor_id" valid:"required~Sponsor is required"`
	Sponsor Sponsor `gorm:"foreignKey:SponsorID" valid:"-"`
	
	ApprovalRequirements    []ApprovalRequirement    `gorm:"foreignKey:ScholarshipID" json:"approval_requirements" valid:"-"`
	ApplicationScholarships []ApplicationScholarship `gorm:"foreignKey:ScholarshipID" json:"application_scholarships" valid:"-"`

	InterviewRounds         []InterviewRound         `gorm:"foreignKey:ScholarshipID" json:"interview_rounds" valid:"-"`
}
