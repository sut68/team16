package test

import (
	"backend/entity"
	"backend/validators"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// unit test for ApprovalDecision

var validAdmin = entity.AdminProfile{
	AdminFirstname: "John",
	AdminLastname:  "Doe",
	Position:       "Manager",
	Email:          "john.doe@example.com",
	Phone:          "1234567890",
	UserID:         1,
}

// case: All fields are valid
func TestApprovalDecisionValidation_AllValid(t *testing.T) {
	g := NewWithT(t)

	decision := entity.ApprovalDecision{
		DecisionAt: time.Now(),
		Decision:   "approve",
		Comment:    "The document is clear and correct.",
		TaskID:     1,
		AdminID:    1,
		Admin:      validAdmin,
	}

	err := validators.ValidateStruct(&decision)

	// Expect no error
	g.Expect(err).To(BeNil())
}

// case: DecisionAt is empty
func TestApprovalDecisionValidation_DecisionAtEmpty(t *testing.T) {
	g := NewWithT(t)

	decision := entity.ApprovalDecision{
		DecisionAt: time.Time{}, // Invalid
		Decision:   "approve",
		Comment:    "This should fail.",
		TaskID:     1,
		AdminID:    1,
		Admin:      validAdmin,
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
		AdminID:    1,
		Admin:      validAdmin,
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
		AdminID:    1,
		Admin:      validAdmin,
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
		AdminID:    1,
		Admin:      validAdmin,
	}

	err := validators.ValidateStruct(&decision)

	// Expect no error because comment is not a required field for validation
	g.Expect(err).To(BeNil())
}

// case: AdminID is zero
func TestApprovalDecisionValidation_AdminIDZero(t *testing.T) {
	g := NewWithT(t)

	decision := entity.ApprovalDecision{
		DecisionAt: time.Now(),
		Decision:   "approve",
		Comment:    "This should fail.",
		TaskID:     1,
		AdminID:    0, // Invalid
		Admin:      validAdmin,
	}

	err := validators.ValidateStruct(&decision)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}
