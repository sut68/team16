package sponsor

import (
	"backend/entity"
	"backend/utils"
	"fmt"

	"gorm.io/gorm"
)

func SeedSponsors(db *gorm.DB) error {
	var industries []entity.SponsorIndustry
	if err := db.Find(&industries).Error; err != nil {
		return fmt.Errorf("failed to fetch industries: %v", err)
	}
	if len(industries) == 0 {
		return fmt.Errorf("industries must be seeded first")
	}

	// Helper สุ่ม industry
	pickIndustry := func(name string) *uint {
		for _, v := range industries {
			if v.Name == name {
				return &v.ID
			}
		}
		// fallback — เอาตัวแรกถ้าไม่เจอ
		id := industries[0].ID
		return &id
	}

	sponsors := []entity.Sponsor {
		{
			CompanyName: "PTT Public Company Limited",
			IndustryID:  pickIndustry("Energy"),
			Website:     utils.StrPtr("https://ptt.example.com"),
			Description: utils.StrPtr("Leading energy and petroleum company in Thailand."),
			Status:      "active",
			Contacts: []entity.SponsorContact{
				{Name: "Nawaporn S.", Email: "nawaporn.ptt@example.com", Phone: "0812345670", Position: utils.StrPtr("Head of Procurement")},
				{Name: "Chaiyapol K.", Email: "chaiyapol.ptt@example.com", Phone: "0891234501", Position: utils.StrPtr("Senior Engineer")},
			},
		},
		{
			CompanyName: "SCG - Siam Cement Group",
			IndustryID:  pickIndustry("Manufacturing"),
			Website:     utils.StrPtr("https://scg.example.com"),
			Description: utils.StrPtr("Top industrial manufacturing and construction solutions provider."),
			Status:      "active",
			Contacts: []entity.SponsorContact{
				{Name: "Somsri L.", Email: "somsri.scg@example.com", Phone: "0861122334", Position: utils.StrPtr("Project Manager")},
				{Name: "Anan P.", Email: "anan.scg@example.com", Phone: "0823344556", Position: utils.StrPtr("Business Development")},
			},
		},
		{
			CompanyName: "AIS - Advanced Info Service",
			IndustryID:  pickIndustry("Telecommunications"),
			Website:     utils.StrPtr("https://ais.example.com"),
			Description: utils.StrPtr("Thailand’s largest telecommunication company."),
			Status:      "active",
			Contacts: []entity.SponsorContact{
				{Name: "Kanya M.", Email: "kanya.ais@example.com", Phone: "0901234567", Position: utils.StrPtr("Sales Engineer")},
				{Name: "Worawit T.", Email: "worawit.ais@example.com", Phone: "0859988776", Position: utils.StrPtr("Network Specialist")},
			},
		},
		{
			CompanyName: "ThaiBev (Thai Beverage)",
			IndustryID:  pickIndustry("Manufacturing"),
			Website:     utils.StrPtr("https://thaibev.example.com"),
			Description: utils.StrPtr("Major beverage manufacturing company in Thailand."),
			Status:      "active",
			Contacts: []entity.SponsorContact{
				{Name: "Pimchanok S.", Email: "pim.thaibev@example.com", Phone: "0819988776", Position: utils.StrPtr("Supply Chain Manager")},
				{Name: "Phumipat R.", Email: "phumipat.thaibev@example.com", Phone: "0894455667", Position: utils.StrPtr("QA Engineer")},
			},
		},
		{
			CompanyName: "CP All Public Company Limited (7-Eleven)",
			IndustryID:  pickIndustry("Industrial Engineering"),
			Website:     utils.StrPtr("https://cpall.example.com"),
			Description: utils.StrPtr("Thailand's largest retail and logistics network operator."),
			Status:      "active",
			Contacts: []entity.SponsorContact{
				{Name: "Suphachai C.", Email: "suphachai.cpall@example.com", Phone: "0867778899", Position: utils.StrPtr("Operations Manager")},
				{Name: "Kanya R.", Email: "kanya.cpall@example.com", Phone: "0829988776", Position: utils.StrPtr("Logistics Coordinator")},
			},
		},
		{
			CompanyName: "Bangkok Airways",
			IndustryID:  pickIndustry("Aerospace"),
			Website:     utils.StrPtr("https://bangkokair.example.com"),
			Description: utils.StrPtr("Regional airline operator in Thailand."),
			Status:      "active",
			Contacts: []entity.SponsorContact{
				{Name: "Narin S.", Email: "narin.bangkokair@example.com", Phone: "0831234987", Position: utils.StrPtr("Operations Engineer")},
				{Name: "May P.", Email: "may.bangkokair@example.com", Phone: "0897766554", Position: utils.StrPtr("Safety Officer")},
			},
		},
		{
			CompanyName: "Delta Electronics (Thailand)",
			IndustryID:  pickIndustry("Electronics"),
			Website:     utils.StrPtr("https://delta.example.com"),
			Description: utils.StrPtr("Leader in power electronics and automation solutions."),
			Status:      "active",
			Contacts: []entity.SponsorContact{
				{Name: "Aroon M.", Email: "aroon.delta@example.com", Phone: "0845566778", Position: utils.StrPtr("Embedded Engineer")},
				{Name: "Ploy T.", Email: "ploy.delta@example.com", Phone: "0873344556", Position: utils.StrPtr("Product Manager")},
			},
		},
		{
			CompanyName: "B.Grimm Power",
			IndustryID:  pickIndustry("Energy"),
			Website:     utils.StrPtr("https://bgrimm.example.com"),
			Description: utils.StrPtr("Power producer and infrastructure company in Thailand."),
			Status:      "active",
			Contacts: []entity.SponsorContact{
				{Name: "Jirayut N.", Email: "jirayut.bgrimm@example.com", Phone: "0903344556", Position: utils.StrPtr("Power Systems Engineer")},
				{Name: "Pimlada K.", Email: "pimlada.bgrimm@example.com", Phone: "0824455667", Position: utils.StrPtr("Project Coordinator")},
			},
		},
		{
			CompanyName: "Kasikornbank (KBank)",
			IndustryID:  pickIndustry("Finance"),
			Website:     utils.StrPtr("https://kbank.example.com"),
			Description: utils.StrPtr("One of Thailand's largest financial institutions."),
			Status:      "active",
			Contacts: []entity.SponsorContact{
				{Name: "Suthida W.", Email: "suthida.kbank@example.com", Phone: "0861230987", Position: utils.StrPtr("Corporate Relations")},
				{Name: "Thiti C.", Email: "thiti.kbank@example.com", Phone: "0899088776", Position: utils.StrPtr("Partnerships Manager")},
			},
		},
		{
			CompanyName: "Toyota Motor Thailand",
			IndustryID:  pickIndustry("Automotive"),
			Website:     utils.StrPtr("https://toyota-th.example.com"),
			Description: utils.StrPtr("Leading automotive manufacturing & engineering company."),
			Status:      "active",
			Contacts: []entity.SponsorContact{
				{Name: "Krit S.", Email: "krit.toyota@example.com", Phone: "0813344556", Position: utils.StrPtr("Mechanical Engineer")},
				{Name: "Natcha R.", Email: "natcha.toyota@example.com", Phone: "0827766554", Position: utils.StrPtr("Supplier Relations")},
			},
		},
	}

	// INSERT
	for _, sponsor := range sponsors {
		var existing entity.Sponsor
		err := db.Where("company_name = ?", sponsor.CompanyName).First(&existing).Error
		
		if err == nil {
			continue
		}
		
		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&sponsor).Error; err != nil {
			return fmt.Errorf("failed to seed sponsor %s: %v", sponsor.CompanyName, err)
		}
	}

	return nil
}