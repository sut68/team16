package test

import (
	"backend/entity"
	"backend/validators"
	"testing"

	. "github.com/onsi/gomega"
)

// unit test for ApprovalTask

// case: All fields are valid
func TestApprovalTaskValidation_AllValid(t *testing.T) {
	g := NewWithT(t)

	task := entity.ApprovalTask{
		Status:     "pending",
		AdminID:    1,
		DocumentID: 1,
	}

	err := validators.ValidateStruct(&task)

	// Expect no error
	g.Expect(err).To(BeNil())
}

// case: Status is empty
func TestApprovalTaskValidation_StatusEmpty(t *testing.T) {
	g := NewWithT(t)

	task := entity.ApprovalTask{
		Status:     "", // Invalid
		AdminID:    1,
		DocumentID: 1,
	}

	err := validators.ValidateStruct(&task)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}

// case: AdminID is zero
func TestApprovalTaskValidation_AdminIDZero(t *testing.T) {
	g := NewWithT(t)

	task := entity.ApprovalTask{
		Status:     "pending",
		AdminID:    0, // Invalid
		DocumentID: 1,
	}

	err := validators.ValidateStruct(&task)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}

// case: DocumentID is zero
func TestApprovalTaskValidation_DocumentIDZero(t *testing.T) {
	g := NewWithT(t)

	task := entity.ApprovalTask{
		Status:     "pending",
		AdminID:    1,
		DocumentID: 0, // Invalid
	}

	err := validators.ValidateStruct(&task)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}
