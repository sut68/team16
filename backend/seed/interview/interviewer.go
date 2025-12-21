package interview

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedInterviewers(db *gorm.DB) {
	interviewers := []entity.Interviewer{
		{InterviewerFirstname: "John", InterviewerLastname: "Doe", Email: "john.doe@example.com"},
		{InterviewerFirstname: "Jane", InterviewerLastname: "Smith", Email: "jane.smith@example.com"},
		{InterviewerFirstname: "Peter", InterviewerLastname: "Jones", Email: "peter.jones@example.com"},
	}

	for _, interviewer := range interviewers {
		db.FirstOrCreate(&interviewer, entity.Interviewer{Email: interviewer.Email})
	}
}
