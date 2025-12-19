package entity

import (
	"time"
	"gorm.io/gorm"
)

type StudentProfile struct {
	gorm.Model

	// Regex เช็ค format มทส. (B/C/M/D ตามด้วยเลข 7 ตัว)
	StudentID string `gorm:"unique;not null" json:"student_id" valid:"required~Student ID is required,matches(^[BCMD][0-9]{7}$)~Invalid Student ID format (e.g. B6630409)"`

	FirstNameTH string `json:"first_name_th" valid:"required~First Name TH is required"`
	LastNameTH  string `json:"last_name_th" valid:"required~Last Name TH is required"`

	// ชื่ออังกฤษต้องเป็นตัวอักษร A-Z เท่านั้น
	FirstNameEN string `json:"first_name_en" valid:"required~First Name EN is required,alpha~First Name EN must be English letters"`
	LastNameEN  string `json:"last_name_en" valid:"required~Last Name EN is required,alpha~Last Name EN must be English letters"`

	// บัตร ปชช. 13 หลัก
	NationalID string `gorm:"unique" json:"national_id" valid:"required~National ID is required,numeric~National ID must be numeric,stringlength(13|13)~National ID must be 13 digits"`

	BirthDate   time.Time `gorm:"type:date" json:"birth_date" valid:"required~Birth Date is required"`
	CurrentYear int       `json:"current_year" valid:"required~Current Year is required,range(1|8)~Current Year must be between 1-8"`
	
	// เกรดเฉลี่ย 0.00 - 4.00
	GPAX float64 `gorm:"type:decimal(3,2)" json:"gpax" valid:"required~GPAX is required,range(0|4)~GPAX must be between 0.00 and 4.00"`

	AdvisorName string `json:"advisor_name" valid:"required~Advisor Name is required"`
	
	// เบอร์โทรศัพท์ 10 หลัก
	Phone string `json:"phone" valid:"required~Phone is required,numeric~Phone must be numeric,stringlength(10|10)~Phone must be 10 digits"`
	
	Email string `json:"email" valid:"required~Email is required,email~Invalid email format"`

	PermanentAddress string `gorm:"type:text" json:"permanent_address" valid:"required~Permanent Address is required"`
	CurrentAddress   string `gorm:"type:text" json:"current_address" valid:"required~Current Address is required"`
	Province         string `json:"province" valid:"required~Province is required"`
	
	// จำนวนพี่น้อง ห้ามติดลบ
	SiblingsCount int `json:"siblings_count" valid:"range(0|20)~Siblings count must be realistic"`

	UserID uint  `json:"user_id" valid:"required~User ID is required"`
	User   *User `gorm:"foreignKey:UserID" json:"user" valid:"-"`

	MajorID uint   `json:"major_id" valid:"required~Major ID is required"`
	Major   *Major `gorm:"foreignKey:MajorID" json:"major" valid:"-"`

	Applications   []Application    `gorm:"foreignKey:StudentProfileID" json:"applications" valid:"-"`
	StudentFavNews []StudentFavNews `gorm:"foreignKey:StudentProfileID" json:"student_fav_news" valid:"-"`

	FamilyInfo *FamilyInfo `gorm:"foreignKey:ProfileID" json:"family_info" valid:"-"`
}