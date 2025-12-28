package scholarship

import (
	"backend/entity"
	"time"

	"gorm.io/gorm"
)

func SeedScholarships(db *gorm.DB) error {
	// Define scholarships to be seeded
	scholarships := []entity.Scholarship{
		// ทุนเต็มจำนวน (Full) - Sponsor 1
		{
			ScholarshipName:     "ทุนเรียนดี ประจำปีการศึกษา 2568",
			Description:         "ทุนการศึกษาสำหรับนักศึกษาที่มีผลการเรียนดีเด่น เกรดเฉลี่ยสะสม 3.50 ขึ้นไป ครอบคลุมค่าเล่าเรียนและค่าครองชีพรายเดือน",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 2, 0).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   1, // Full
			SemasterID:          1,
			SponsorID:           1,
		},
		// ทุนบางส่วน (Partial) - Sponsor 1
		{
			ScholarshipName:     "ทุนช่วยเหลือนักศึกษาขาดแคลนทุนทรัพย์",
			Description:         "ทุนสนับสนุนการศึกษาสำหรับนักศึกษาที่ขาดแคลนทุนทรัพย์ ช่วยเหลือค่าเล่าเรียนบางส่วน พิจารณาจากฐานะครอบครัว",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 1, 15).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   2, // Partial
			SemasterID:          1,
			SponsorID:           1,
		},
		// ทุนเต็มจำนวน (Full) - Sponsor 2
		{
			ScholarshipName:     "ทุนนักกีฬาดีเด่น",
			Description:         "ทุนการศึกษาสำหรับนักศึกษาที่มีความสามารถด้านกีฬาระดับชาติ หรือได้รับเหรียญรางวัลจากการแข่งขัน ครอบคลุมค่าเล่าเรียนทั้งหมด",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 3, 0).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   1, // Full
			SemasterID:          1,
			SponsorID:           2,
		},
		// ทุนบางส่วน (Partial) - Sponsor 2
		{
			ScholarshipName:     "ทุนส่งเสริมศิลปวัฒนธรรม",
			Description:         "ทุนสนับสนุนนักศึกษาที่มีความสามารถด้านศิลปะและวัฒนธรรมไทย เช่น นาฏศิลป์ ดนตรีไทย จิตรกรรม รับทุนค่าเล่าเรียน 50%",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 2, 0).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   2, // Partial
			SemasterID:          1,
			SponsorID:           2,
		},
		// ทุนเต็มจำนวน (Full) - Sponsor 3
		{
			ScholarshipName:     "ทุนพัฒนาบุคลากรด้านเทคโนโลยี",
			Description:         "ทุนการศึกษาเต็มจำนวนสำหรับนักศึกษาสาขาวิศวกรรมคอมพิวเตอร์ วิทยาการคอมพิวเตอร์ หรือสาขาที่เกี่ยวข้องกับ IT พร้อมโอกาสฝึกงาน",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   1, // Full
			SemasterID:          1,
			SponsorID:           3,
		},
		// ทุนปิดรับสมัครแล้ว (Closed) - เพื่อทดสอบ
		{
			ScholarshipName:     "ทุนวิจัยระดับบัณฑิตศึกษา (รอบที่ 1)",
			Description:         "ทุนสนับสนุนการวิจัยสำหรับนักศึกษาระดับปริญญาโทและเอก ครอบคลุมค่าใช้จ่ายในการทำวิจัยและค่าเล่าเรียน",
			OpenDate:            time.Now().AddDate(0, -3, 0).Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, -1, 0).Format("2006-01-02"),
			StatusscholarshipID: 2, // Closed
			TypescholarshipID:   1, // Full
			SemasterID:          1,
			SponsorID:           1,
		},
		// ทุนบางส่วน (Partial) - Sponsor 3
		{
			ScholarshipName:     "ทุนจิตอาสาพัฒนาชุมชน",
			Description:         "ทุนสำหรับนักศึกษาที่มีจิตอาสา ทำกิจกรรมเพื่อสังคมและชุมชน มีชั่วโมงจิตอาสาไม่น้อยกว่า 50 ชั่วโมง/ภาคเรียน",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 1, 20).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   2, // Partial
			SemasterID:          1,
			SponsorID:           3,
		},
		// ทุนเต็มจำนวน (Full) - Sponsor 1
		{
			ScholarshipName:     "ทุนเฉลิมพระเกียรติ",
			Description:         "ทุนพระราชทานเพื่อส่งเสริมการศึกษาแก่นักศึกษาที่มีความประพฤติดี มีจิตสาธารณะ และมีผลการเรียนดี เกรดเฉลี่ย 3.25 ขึ้นไป",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 2, 15).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   1, // Full
			SemasterID:          1,
			SponsorID:           1,
		},

		// Sponsor 4: ThaiBev
		{
			ScholarshipName:     "ทุนไทยเบฟส่งเสริมการศึกษา",
			Description:         "ทุนการศึกษาจากบริษัท ไทยเบฟเวอเรจ จำกัด สำหรับนักศึกษาสาขาบริหารธุรกิจ การตลาด หรือสาขาที่เกี่ยวข้อง ครอบคลุมค่าเล่าเรียนเต็มจำนวน",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 2, 0).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   1, // Full
			SemasterID:          1,
			SponsorID:           4,
		},

		// Sponsor 5: CP All
		{
			ScholarshipName:     "ทุน CP ALL เพื่อพัฒนาบุคลากรค้าปลีก",
			Description:         "ทุนการศึกษาสำหรับนักศึกษาที่สนใจธุรกิจค้าปลีก โลจิสติกส์ และการจัดการซัพพลายเชน พร้อมโอกาสฝึกงานและทำงานกับ 7-Eleven",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 1, 10).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   2, // Partial
			SemasterID:          1,
			SponsorID:           5,
		},

		// Sponsor 6: Bangkok Airways
		{
			ScholarshipName:     "ทุนบางกอกแอร์เวย์สสู่ท้องฟ้า",
			Description:         "ทุนการศึกษาสำหรับนักศึกษาสาขาวิศวกรรมการบิน การจัดการการบิน หรือสาขาที่เกี่ยวข้อง สนับสนุนค่าเล่าเรียนเต็มจำนวนพร้อมการฝึกภาคปฏิบัติ",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 3, 0).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   1, // Full
			SemasterID:          1,
			SponsorID:           6,
		},

		// Sponsor 7: Delta Electronics
		{
			ScholarshipName:     "ทุน Delta เพื่ออนาคตอิเล็กทรอนิกส์",
			Description:         "ทุนการศึกษาจาก Delta Electronics สำหรับนักศึกษาสาขาวิศวกรรมไฟฟ้า อิเล็กทรอนิกส์ หรือระบบฝังตัว ครอบคลุมค่าเล่าเรียน 75%",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 2, 0).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   2, // Partial
			SemasterID:          1,
			SponsorID:           7,
		},

		// Sponsor 8: B.Grimm Power
		{
			ScholarshipName:     "ทุน B.Grimm พลังงานสะอาดเพื่ออนาคต",
			Description:         "ทุนการศึกษาสำหรับนักศึกษาที่สนใจด้านพลังงานหมุนเวียน พลังงานแสงอาทิตย์ และวิศวกรรมพลังงาน ครอบคลุมค่าเล่าเรียนเต็มจำนวน",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 2, 20).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   1, // Full
			SemasterID:          1,
			SponsorID:           8,
		},

		// Sponsor 9: Kasikornbank (KBank)
		{
			ScholarshipName:     "ทุนกสิกรไทยสร้างอนาคต",
			Description:         "ทุนการศึกษาจากธนาคารกสิกรไทย สำหรับนักศึกษาสาขาการเงิน บัญชี เศรษฐศาสตร์ หรือ FinTech พร้อมโปรแกรมพัฒนาทักษะการเงิน",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 1, 25).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   1, // Full
			SemasterID:          1,
			SponsorID:           9,
		},

		// Sponsor 10: Toyota
		{
			ScholarshipName:     "ทุน Toyota วิศวกรรมยานยนต์แห่งอนาคต",
			Description:         "ทุนการศึกษาจาก Toyota Motor Thailand สำหรับนักศึกษาสาขาวิศวกรรมยานยนต์ เครื่องกล หรือการผลิต ครอบคลุมค่าเล่าเรียนและค่าอุปกรณ์การเรียน",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 2, 10).Format("2006-01-02"),
			StatusscholarshipID: 1, // Open
			TypescholarshipID:   1, // Full
			SemasterID:          1,
			SponsorID:           10,
		},
	}

	for _, s := range scholarships {
		// Use FirstOrCreate to prevent duplicates on subsequent runs
		if err := db.Where("scholarship_name = ?", s.ScholarshipName).FirstOrCreate(&s).Error; err != nil {
			return err
		}
	}

	return nil
}
