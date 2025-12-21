package location

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedLocations(db *gorm.DB) {
	locations := []entity.Location{
		{Name: "อาคารเรียนรวม 1 - ห้อง 101", Building: "อาคารเรียนรวม 1", Room: "101", Floor: 1, Description: ""},
		{Name: "อาคารวิศวกรรมศาสตร์ - ห้องประชุม 2", Building: "อาคารวิศวกรรมศาสตร์", Room: "ประชุม 2", Floor: 2, Description: ""},
		{Name: "Online Meeting", Building: "Online", Room: "Online", Floor: 0, Description: "ช่องทางออนไลน์"},
	}

	for _, loc := range locations {
		db.FirstOrCreate(&loc, entity.Location{Name: loc.Name})
	}
}
