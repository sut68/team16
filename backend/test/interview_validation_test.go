package test

import (
	"backend/entity"
	"backend/validators"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestInterviewRoundFullValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	interviewModeID := uint(1)
	locationID := uint(1)

	round := entity.InterviewRound{
		Name:            "Test Round",
		Description:     "Test Description",
		StartDateTime:   time.Now().Add(24 * time.Hour),
		EndDateTime:     time.Now().Add(48 * time.Hour),
		SlotDuration:    30,
		ScholarshipID:   1,
		AdminProfileID:  1,
		InterviewModeID: &interviewModeID,
		LocationID:      &locationID,
		MeetingLink:     "http://example.com",
	}

	err := validators.ValidateStruct(round)
	g.Expect(err).To(BeNil())
}

func TestInterviewRoundNameValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	interviewModeID := uint(1)

	t.Run("Name is required", func(t *testing.T) {
		round := entity.InterviewRound{
			Name:            "", // Invalid
			Description:     "Test Description",
			StartDateTime:   time.Now().Add(24 * time.Hour),
			EndDateTime:     time.Now().Add(48 * time.Hour),
			SlotDuration:    30,
			ScholarshipID:   1,
			AdminProfileID:  1,
			InterviewModeID: &interviewModeID,
		}

		err := validators.ValidateStruct(round)
		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Error()).To(Equal("Name is required"))
	})
}

func TestInterviewRoundSlotDurationValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	interviewModeID := uint(1)

	t.Run("SlotDuration cannot be zero", func(t *testing.T) {
		round := entity.InterviewRound{
			Name:            "Test Round",
			Description:     "Test Description",
			StartDateTime:   time.Now().Add(24 * time.Hour),
			EndDateTime:     time.Now().Add(48 * time.Hour),
			SlotDuration:    0, // Invalid
			ScholarshipID:   1,
			AdminProfileID:  1,
			InterviewModeID: &interviewModeID,
		}

		err := validators.ValidateStruct(round)
		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Error()).To(Equal("slot_duration: non zero value required"))
	})
}

func TestInterviewRoundForeignKeyValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	interviewModeID := uint(1)

	t.Run("ScholarshipID is required", func(t *testing.T) {
		round := entity.InterviewRound{
			Name:            "Test Round",
			Description:     "Test Description",
			StartDateTime:   time.Now().Add(24 * time.Hour),
			EndDateTime:     time.Now().Add(48 * time.Hour),
			SlotDuration:    30,
			ScholarshipID:   0, // Invalid
			AdminProfileID:  1,
			InterviewModeID: &interviewModeID,
		}

		err := validators.ValidateStruct(round)
		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Error()).To(Equal("ScholarshipID is required"))
	})
}
