package approval

import (
	"backend/entity"
	"errors"

	"gorm.io/gorm"
)

func SeedApprovalRequirements(db *gorm.DB) error {
	if err := db.First(&entity.ApprovalRequirement{}).Error; err == gorm.ErrRecordNotFound {
		// Find scholarship
		var scholarship entity.Scholarship
		if err := db.Where("scholarship_name = ?", "Beasiswa Anak Bangsa").First(&scholarship).Error; err != nil {
			return errors.New("scholarship 'Beasiswa Anak Bangsa' not found for seeding ApprovalRequirements")
		}

		// Find master requirements
		var req1, req2 entity.Requirement
		if err := db.Where("name = ?", "Transkrip nilai semester terakhir").First(&req1).Error; err != nil {
			return errors.New("requirement 'Transkrip nilai semester terakhir' not found")
		}
		if err := db.Where("name = ?", "Surat rekomendasi dari dekan fakultas").First(&req2).Error; err != nil {
			return errors.New("requirement 'Surat rekomendasi dari dekan fakultas' not found")
		}

		requirements := []entity.ApprovalRequirement{
			{
				ScholarshipID: scholarship.ID,
				RequirementID: req1.ID,
			},
			{
				ScholarshipID: scholarship.ID,
				RequirementID: req2.ID,
			},
		}

		if err := db.Create(&requirements).Error; err != nil {
			return err
		}
	}
	return nil
}

// func SeedApprovalRequirements(db *gorm.DB) error {
//     if err := db.First(&entity.ApprovalRequirement{}).Error; err == gorm.ErrRecordNotFound {
        
//         // ----------------- ส่วนของทุนที่ 1 (เดิม) -----------------
//         var scholarship1 entity.Scholarship
//         if err := db.Where("scholarship_name = ?", "Beasiswa Anak Bangsa").First(&scholarship1).Error; err != nil {
//             return errors.New("scholarship 1 not found")
//         }

//         // ค้นหา Master Requirements
//         var req1, req2 entity.Requirement
//         if err := db.Where("name = ?", "Transkrip nilai semester terakhir").First(&req1).Error; err != nil { return err }
//         if err := db.Where("name = ?", "Surat rekomendasi dari dekan fakultas").First(&req2).Error; err != nil { return err }

//         // สร้าง Requirement ให้ทุนที่ 1
//         reqs1 := []entity.ApprovalRequirement{
//             { ScholarshipID: scholarship1.ID, RequirementID: req1.ID },
//             { ScholarshipID: scholarship1.ID, RequirementID: req2.ID },
//         }
//         if err := db.Create(&reqs1).Error; err != nil { return err }


//         // ----------------- 🔥 ส่วนที่ต้องเพิ่มสำหรับทุนที่ 2 -----------------
//         var scholarship2 entity.Scholarship
//         // ค้นหาทุน ID 2 (หรือค้นหาตามชื่อที่คุณตั้งใน SeedScholarships)
//         if err := db.First(&scholarship2, 2).Error; err == nil { 
            
//             reqs2 := []entity.ApprovalRequirement{
//                 {
//                     ScholarshipID: scholarship2.ID, // ผูกกับทุนที่ 2
//                     RequirementID: req1.ID,         // ใช้เงื่อนไขเดียวกัน (Transcript)
//                 },
//                 {
//                     ScholarshipID: scholarship2.ID,
//                     RequirementID: req2.ID,         // ใช้เงื่อนไขเดียวกัน (Surat)
//                 },
//             }
            
//             if err := db.Create(&reqs2).Error; err != nil {
//                 return err
//             }
//         }
//         // -------------------------------------------------------------
//     }
//     return nil
// }