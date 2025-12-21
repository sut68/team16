package test

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"backend/entity"
	"backend/validators"

)

// unit test for ApprovalDecision

// case: All fields are valid
func TestApprovalDecisionValidation_AllValid(t *testing.T) {
	g := NewWithT(t)

	decision := entity.ApprovalDecision{
		DecisionAt: time.Now(),
		Decision:   "approve",
		Comment:    "The document is clear and correct.",
		TaskID:     1,
	}

	err := validators.ValidateStruct(&decision)

	// Expect no error
	g.Expect(err).To(BeNil())
}

// case: DecisionAt is empty
func TestApprovalDecisionValidation_DecisionAtEmpty(t *testing.T) {
	g := NewWithT(t)

	decision := entity.ApprovalDecision{
		DecisionAt: time.Now(), // Invalid
		Decision:   "approve",
		Comment:    "This should fail.",
		TaskID:     1,
	}

	err := validators.ValidateStruct(&decision)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}

// case: Decision is empty
func TestApprovalDecisionValidation_DecisionEmpty(t *testing.T) {
	g := NewWithT(t)

	decision := entity.ApprovalDecision{
		DecisionAt: time.Now(),
		Decision:   "", // Invalid
		Comment:    "This should fail.",
		TaskID:     1,
	}

	err := validators.ValidateStruct(&decision)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}

// case: TaskID is zero
func TestApprovalDecisionValidation_TaskIDZero(t *testing.T) {
	g := NewWithT(t)

	decision := entity.ApprovalDecision{
		DecisionAt: time.Now(),
		Decision:   "reject",
		Comment:    "This should fail.",
		TaskID:     0, // Invalid
	}

	err := validators.ValidateStruct(&decision)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}

// case: Comment is empty (which is allowed)
func TestApprovalDecisionValidation_CommentEmpty(t *testing.T) {
	g := NewWithT(t)

	decision := entity.ApprovalDecision{
		DecisionAt: time.Now(),
		Decision:   "approve",
		Comment:    "", // Allowed
		TaskID:     1,
	}

	err := validators.ValidateStruct(&decision)

	// Expect no error because comment is not a required field for validation
	g.Expect(err).To(BeNil())
}
