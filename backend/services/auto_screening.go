package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"backend/config"
	"backend/entity"
	"backend/ws"
)

// AutoScreenResult represents the result of auto screening for one application
type AutoScreenResult struct {
	ScreeningID      uint     `json:"screening_id"`
	ApplicationID    uint     `json:"application_id"`
	ScholarshipID    uint     `json:"scholarship_id"`
	ScholarshipName  string   `json:"scholarship_name"`
	StudentName      string   `json:"student_name"`
	AutoApproved     bool     `json:"auto_approved"` // true = auto approved, false = needs manual review
	PassedCriteria   int      `json:"passed_criteria"`
	TotalCriteria    int      `json:"total_criteria"`
	FailedCriteria   []string `json:"failed_criteria,omitempty"`
	ProcessedByAdmin string   `json:"processed_by_admin,omitempty"`
}

// CriteriaCheckResult represents result of checking one criterion
type CriteriaCheckResult struct {
	CriteriaName  string  `json:"criteria_name"`
	Operator      string  `json:"operator"`
	RequiredValue float64 `json:"required_value"`
	StudentValue  float64 `json:"student_value"`
	Passed        bool    `json:"passed"`
	FailureReason string  `json:"failure_reason,omitempty"`
}

const (
	STATUS_PENDING  = 1
	STATUS_APPROVED = 2
	STATUS_REJECTED = 3
)

// AutoScreenSingle performs auto screening for a single screening record
// Returns the result and updates the database if all criteria pass
func AutoScreenSingle(screeningID uint, adminProfileID uint, adminName string) (*AutoScreenResult, error) {
	var screening entity.Screening

	// Load screening with all related data
	if err := config.DB.
		Preload("ApplicationScholarship.Application.StudentProfile.FamilyInfo").
		Preload("ApplicationScholarship.Scholarship").
		First(&screening, screeningID).Error; err != nil {
		return nil, fmt.Errorf("screening not found: %w", err)
	}

	// Skip if already processed
	if screening.StatusScreeningID != STATUS_PENDING {
		return nil, fmt.Errorf("screening already processed (status: %d)", screening.StatusScreeningID)
	}

	appSch := screening.ApplicationScholarship
	scholarship := appSch.Scholarship
	studentPtr := appSch.Application.StudentProfile

	// Handle nil StudentProfile
	if studentPtr == nil {
		return nil, fmt.Errorf("student profile not found")
	}
	student := *studentPtr

	// Handle nil FamilyInfo
	var family entity.FamilyInfo
	if student.FamilyInfo != nil {
		family = *student.FamilyInfo
	}

	// Load scholarship features (criteria)
	var features []entity.Featurescholarship
	if err := config.DB.
		Preload("Typefeature").
		Where("scholarship_id = ?", scholarship.ID).
		Find(&features).Error; err != nil {
		return nil, fmt.Errorf("failed to load features: %w", err)
	}

	// Check all criteria
	checkResults := checkAllCriteria(features, student, family)

	// Count passed/failed
	passedCount := 0
	var failedCriteria []string
	for _, result := range checkResults {
		if result.Passed {
			passedCount++
		} else {
			failedCriteria = append(failedCriteria, result.FailureReason)
		}
	}

	totalCriteria := len(checkResults)
	allPassed := passedCount == totalCriteria && totalCriteria > 0

	result := &AutoScreenResult{
		ScreeningID:      screeningID,
		ApplicationID:    appSch.ApplicationID,
		ScholarshipID:    scholarship.ID,
		ScholarshipName:  scholarship.ScholarshipName,
		StudentName:      fmt.Sprintf("%s %s", student.FirstNameTH, student.LastNameTH),
		AutoApproved:     allPassed,
		PassedCriteria:   passedCount,
		TotalCriteria:    totalCriteria,
		FailedCriteria:   failedCriteria,
		ProcessedByAdmin: adminName,
	}

	// Only auto-approve if ALL criteria pass
	if allPassed {
		tx := config.DB.Begin()

		// Update screening status to APPROVED
		if err := tx.Model(&screening).Updates(map[string]interface{}{
			"status_screening_id": STATUS_APPROVED,
			"admin_profile_id":    adminProfileID,
			"rejection_reason":    nil,
		}).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update screening: %w", err)
		}

		// Update ApplicationScholarship status (ready for document upload)
		if err := tx.Model(&entity.ApplicationScholarship{}).
			Where("id = ?", appSch.ID).
			Update("status", "").Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update application status: %w", err)
		}

		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Errorf("failed to commit: %w", err)
		}

		log.Printf("✅ Auto-approved screening #%d for %s", screeningID, result.StudentName)
	} else {
		// Don't update - leave for manual review
		log.Printf("⏳ Screening #%d for %s needs manual review (passed %d/%d)",
			screeningID, result.StudentName, passedCount, totalCriteria)
	}

	// Broadcast result via WebSocket
	if ws.ScreeningHubInstance != nil {
		ws.ScreeningHubInstance.BroadcastScreeningResult(ws.ScreeningResult{
			ScreeningID:     screeningID,
			ApplicationID:   appSch.ApplicationID,
			ScholarshipID:   scholarship.ID,
			ScholarshipName: scholarship.ScholarshipName,
			StudentName:     result.StudentName,
			Passed:          allPassed,
			PassedCriteria:  passedCount,
			TotalCriteria:   totalCriteria,
			FailedReasons:   failedCriteria,
			ProcessedBy:     adminName,
		})
	}

	return result, nil
}

// AutoScreenBatch performs auto screening for all pending screenings of a scholarship
func AutoScreenBatch(scholarshipID uint, adminProfileID uint, adminName string) ([]AutoScreenResult, error) {
	var screenings []entity.Screening

	// Find all pending screenings for this scholarship
	if err := config.DB.
		Joins("JOIN application_scholarships ON application_scholarships.id = screenings.application_scholarship_id").
		Where("application_scholarships.scholarship_id = ? AND screenings.status_screening_id = ?", scholarshipID, STATUS_PENDING).
		Find(&screenings).Error; err != nil {
		return nil, fmt.Errorf("failed to find screenings: %w", err)
	}

	if len(screenings) == 0 {
		return []AutoScreenResult{}, nil
	}

	// Get scholarship name for broadcasting
	var scholarship entity.Scholarship
	config.DB.First(&scholarship, scholarshipID)

	var results []AutoScreenResult
	processed := 0
	passed := 0
	failed := 0

	for _, screening := range screenings {
		result, err := AutoScreenSingle(screening.ID, adminProfileID, adminName)
		if err != nil {
			log.Printf("Error processing screening #%d: %v", screening.ID, err)
			continue
		}

		results = append(results, *result)
		processed++

		if result.AutoApproved {
			passed++
		} else {
			failed++
		}

		// Broadcast progress
		if ws.ScreeningHubInstance != nil {
			ws.ScreeningHubInstance.BroadcastProgress(ws.BatchProgress{
				ScholarshipID:   scholarshipID,
				ScholarshipName: scholarship.ScholarshipName,
				Total:           len(screenings),
				Processed:       processed,
				Passed:          passed,
				Failed:          failed,
			})
		}
	}

	// Broadcast completion
	if ws.ScreeningHubInstance != nil {
		ws.ScreeningHubInstance.BroadcastBatchComplete(scholarshipID, scholarship.ScholarshipName, len(screenings), passed, failed)
	}

	return results, nil
}

// checkAllCriteria checks all scholarship criteria against student data
func checkAllCriteria(features []entity.Featurescholarship, student entity.StudentProfile, family entity.FamilyInfo) []CriteriaCheckResult {
	var results []CriteriaCheckResult

	for _, feature := range features {
		// Get type feature name (Typefeaturename is the correct field name)
		typeName := feature.Typefeature.Typefeaturename

		criteriaName := typeName
		if criteriaName == "" {
			criteriaName = fmt.Sprintf("เกณฑ์ #%d", feature.ID)
		}

		fullText := strings.ToLower(criteriaName)
		operator := feature.Operator
		if operator == "" {
			operator = ">="
		}

		requiredValue, _ := strconv.ParseFloat(feature.Value, 64)

		var studentValue float64
		var unit string

		// Determine which value to check based on criteria name
		switch {
		case strings.Contains(fullText, "เกรด") || strings.Contains(fullText, "gpax"):
			// GPAX is float64 in entity
			studentValue = student.GPAX
			unit = ""

		case strings.Contains(fullText, "รายได้") || strings.Contains(fullText, "income"):
			// Calculate family income - FatherIncome, MotherIncome, GuardianIncome are float64
			totalIncome := family.FatherIncome + family.MotherIncome
			if family.GuardianIsParent == "other" || family.GuardianIsParent == "" {
				// Only add guardian income if they're not already counted as father/mother
				if family.GuardianName != "" && family.GuardianName != family.FatherName && family.GuardianName != family.MotherName {
					totalIncome += family.GuardianIncome
				}
			}

			if strings.Contains(fullText, "ต่อคน") || strings.Contains(fullText, "เฉลี่ย") || strings.Contains(fullText, "สมาชิก") {
				// Per capita income
				memberCount := 1 + student.SiblingsCount
				if family.FatherIncome > 0 || family.FatherName != "" {
					memberCount++
				}
				if family.MotherIncome > 0 || family.MotherName != "" {
					memberCount++
				}
				if (family.GuardianIsParent == "other" || family.GuardianIsParent == "") &&
					(family.GuardianIncome > 0 || family.GuardianName != "") {
					if family.GuardianName != family.FatherName && family.GuardianName != family.MotherName {
						memberCount++
					}
				}
				if memberCount == 0 {
					memberCount = 1
				}
				studentValue = totalIncome / float64(memberCount)
				unit = "บาท/คน"
			} else {
				studentValue = totalIncome
				unit = "บาท"
			}

		case strings.Contains(fullText, "ชั้นปี") || strings.Contains(fullText, "ระยะเวลา"):
			studentValue = float64(student.CurrentYear)
			unit = "ปี"

		case strings.Contains(fullText, "พี่น้อง"):
			studentValue = float64(student.SiblingsCount)
			unit = "คน"

		default:
			// Unknown criteria - skip
			continue
		}

		// Compare values
		passed := compareValues(studentValue, operator, requiredValue)

		result := CriteriaCheckResult{
			CriteriaName:  criteriaName,
			Operator:      operator,
			RequiredValue: requiredValue,
			StudentValue:  studentValue,
			Passed:        passed,
		}

		if !passed {
			result.FailureReason = formatFailureReason(criteriaName, operator, requiredValue, studentValue, unit)
		}

		results = append(results, result)
	}

	return results
}

// compareValues compares student value against required value using operator
func compareValues(studentValue float64, operator string, requiredValue float64) bool {
	switch operator {
	case ">=":
		return studentValue >= requiredValue
	case "<=":
		return studentValue <= requiredValue
	case ">":
		return studentValue > requiredValue
	case "<":
		return studentValue < requiredValue
	case "=", "==":
		return studentValue == requiredValue
	default:
		return studentValue >= requiredValue
	}
}

// formatFailureReason creates a human-readable failure reason
func formatFailureReason(criteriaName, operator string, required, actual float64, unit string) string {
	opText := map[string]string{
		">=": "ไม่ต่ำกว่า",
		"<=": "ไม่เกิน",
		">":  "มากกว่า",
		"<":  "น้อยกว่า",
		"=":  "เท่ากับ",
		"==": "เท่ากับ",
	}

	op := opText[operator]
	if op == "" {
		op = operator
	}

	return fmt.Sprintf("%s: ค่าที่ได้ %.2f %s (เกณฑ์: %s %.2f %s)",
		criteriaName, actual, unit, op, required, unit)
}
