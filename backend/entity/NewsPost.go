package entity

import "gorm.io/gorm"

type NewsPost struct {
    gorm.Model
	
    Title string `gorm:"not null" json:"title" valid:"required~You must provide a title"`

    FilePath string `gorm:"not null" json:"file_path" valid:"required~File path is required"`
    
    PostDetail string `gorm:"not null" json:"post_detail" valid:"required~You must provide post details"`

    AdminID uint         `json:"admin_id" gorm:"not null" valid:"required"`
    Admin   AdminProfile `json:"admin_profile" gorm:"foreignKey:AdminID" valid:"-"`

    ScholarshipID uint        `json:"scholarship_id" gorm:"not null" valid:"required"`
    Scholarship   Scholarship `json:"scholarship" gorm:"foreignKey:ScholarshipID" valid:"-"`

    StatusNewsID uint       `json:"status_news_id" gorm:"not null" valid:"required,in(1|2|3|4|5)~Invalid status"`
    StatusNews   StatusNews `json:"status_news" gorm:"foreignKey:StatusNewsID" valid:"-"`
}

/*
valid:"required" → ห้ามว่าง

optional → ไม่ส่งมาก็ได้ 

stringlength(a|b) → ความยาว

in(1|2|3) → ค่า enum

logic ซับซ้อน → เขียนใน BeforeCreate
*/