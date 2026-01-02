package seed

import (
	"log"

	"gorm.io/gorm"

	"backend/seed/approval"
	"backend/seed/evaluation"
	"backend/seed/interview"
	"backend/seed/location"
	"backend/seed/news"
	"backend/seed/scholarship"
	"backend/seed/screening"
	"backend/seed/semaster"
	"backend/seed/sponsor"
	"backend/seed/user"
)

func SeedAll(db *gorm.DB) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	location.SeedLocations(db)
	log.Println("Location seed completed")

	interview.SeedInterviewers(db)
	log.Println("Interviewer seed completed")

	interview.SeedInterviewModes(db)
	log.Println("Interview mode seed completed")

	if err := sponsor.SeedIndustries(db); err != nil {
		return err
	}
	log.Println("Sponsor industries seed completed")

	if err := sponsor.SeedSponsors(db); err != nil {
		return err
	}
	log.Println("Sponsors seed completed")

	if err := user.SeedRoles(db); err != nil {
		return err
	}
	log.Println("Roles seed completed")

	if err := user.SeedUsers(db); err != nil {
		return err
	}
	log.Println("Users seed completed")

	if err := user.SeedMajors(db); err != nil {
		return err
	}
	log.Println("Majors seed completed")

	if err := user.SeedStudentProfiles(db); err != nil {
		return err
	}
	log.Println("Student profiles seed completed")

	if err := user.SeedAdminProfiles(db); err != nil {
		return err
	}
	log.Println("Admin profiles seed completed")

	if err := user.SeedFamilyInfos(db); err != nil {
		return err
	}
	log.Println("Family infos seed completed")

	if err := semaster.CreateSemasters(db); err != nil {
		return err
	}
	log.Println("Semesters seed completed")

	if err := scholarship.SeedStatusScholarships(db); err != nil {
		return err
	}
	log.Println("Status scholarships seed completed")

	if err := scholarship.SeedTypeScholarships(db); err != nil {
		return err
	}
	log.Println("Type scholarships seed completed")

	if err := scholarship.SeedTypeFeatures(db); err != nil {
		return err
	}
	log.Println("Type features seed completed")

	if err := scholarship.SeedScholarships(db); err != nil {
		return err
	}
	log.Println("Scholarships seed completed")

	if err := scholarship.SeedFeatureScholarships(db); err != nil {
		return err
	}
	log.Println("Feature scholarships seed completed")

	if err := scholarship.SeedRequirements(db); err != nil {
		return err
	}
	log.Println("Requirements seed completed")

	if err := approval.SeedApplications(db); err != nil {
		return err
	}
	log.Println("Applications seed completed")

	if err := approval.SeedApplicationScholarships(db); err != nil {
		return err
	}
	log.Println("Application scholarships seed completed")

	if err := approval.SeedApplicationDocuments(db); err != nil {
		return err
	}
	log.Println("Application documents seed completed")

	if err := approval.SeedApprovalRequirements(db); err != nil {
		return err
	}
	log.Println("Approval requirements seed completed")

	// --- New/Corrected Seeders ---
	if err := approval.SeedApprovalTasks(db); err != nil {
		return err
	}
	log.Println("Approval tasks seed completed")

	if err := approval.SeedApprovalDecisions(db); err != nil {
		return err
	}
	log.Println("Approval decisions seed completed")

	if err := screening.SeedStatusScreenings(db); err != nil {
		return err
	}
	log.Println("Status screenings seed completed")

	if err := news.SeedStatusNews(db); err != nil {
		return err
	}
	log.Println("Status news seed completed")

	if err := evaluation.SeedEvaluationCriteria(db); err != nil {
		return err
	}
	log.Println("Evaluation criteria seed completed")

	if err := interview.SeedInterviewRounds(db); err != nil {
		return err
	}
	log.Println("Interview rounds and slots seed completed")

	if err := interview.SeedInterviewBookings(db); err != nil {
		return err
	}
	log.Println("Interview bookings seed completed")

	if err := evaluation.SeedInterviewRoundCriteria(db); err != nil {
		return err
	}
	log.Println("Interview round criteria seed completed")

	if err := evaluation.SeedEvaluations(db); err != nil {
		return err
	}
	log.Println("Evaluations seed completed")

	if err := evaluation.SeedEvaluationScores(db); err != nil {
		return err
	}
	log.Println("Evaluation scores seed completed")

	return nil
}
