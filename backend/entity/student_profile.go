package entity

import "gorm.io/gorm"
import "time"

type StudentProfile struct {
	gorm.Model
	StudentID        string    `gorm:"unique;not null" json:"student_id"`
	FirstNameTH      string    `json:"first_name_th"`
	LastNameTH       string    `json:"last_name_th"`
	FirstNameEN      string    `json:"first_name_en"`
	LastNameEN       string    `json:"last_name_en"`
	NationalID       string    `gorm:"unique" json:"national_id"`
	BirthDate        time.Time `gorm:"type:date" json:"birth_date"`
	CurrentYear      int       `json:"current_year"`
	GPAX             float64   `gorm:"type:decimal(3,2)" json:"gpax"`
	AdvisorName      string    `json:"advisor_name"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email"`
	PermanentAddress string    `gorm:"type:text" json:"permanent_address"`
	CurrentAddress   string    `gorm:"type:text" json:"current_address"`
	Province         string    `json:"province"`
	SiblingsCount    int       `json:"siblings_count"`

	UserID uint `json:"user_id"`
	User   *User `gorm:"foreignKey:UserID" json:"user"`

	MajorID uint  `json:"major_id"`
	Major   *Major `gorm:"foreignKey:MajorID" json:"major"`

	Applications []Application `gorm:"foreignKey:StudentProfileID" json:"applications"`
	StudentFavNews []StudentFavNews `gorm:"foreignKey:StudentProfileID" json:"student_fav_news"`
}
