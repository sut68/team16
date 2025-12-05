package entity

import "gorm.io/gorm"

type NewsPost struct {
	gorm.Model

	Title string `gorm:"not null" json:"title"`

	FilePath   string `gorm:"not null" json:"file_path"`
	PostDetail string `gorm:"not null" json:"post_detail"`

	// ---------- AdminProfile ----------
	AdminID uint         `json:"admin_id" gorm:"not null"`
	Admin   AdminProfile `json:"admin_profile" gorm:"foreignKey:AdminID"`

	// ---------- Scholarship ----------
	ScholarshipID uint        `json:"scholarship_id" gorm:"not null"`
	Scholarship   Scholarship `json:"scholarship" gorm:"foreignKey:ScholarshipID"`

	// ---------- Status NewsID ----------
	StatusNewsID uint       `json:"status_news_id" gorm:"not null"`
	StatusNews   StatusNews `json:"status_news" gorm:"foreignKey:StatusNewsID"`

}
