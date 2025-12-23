package test

import (
	"backend/entity"
	"backend/validators"
	"testing"

	. "github.com/onsi/gomega"
)

// unit test for ApprovalTask

var validDoc = entity.ApplicationDocument{
	FileName:                 "document.pdf",
	UploadedBy:               "Student",
	ApplicationScholarshipID: 1,
}

// case: All fields are valid
func TestApprovalTaskValidation_AllValid(t *testing.T) {
	g := NewWithT(t)

	task := entity.ApprovalTask{
		Status:              "pending",
		DocumentID:          1,
		ApplicationDocument: validDoc,
	}

	err := validators.ValidateStruct(&task)

	// Expect no error
	g.Expect(err).To(BeNil())
}

// case: Status is empty
func TestApprovalTaskValidation_StatusEmpty(t *testing.T) {
	g := NewWithT(t)

	task := entity.ApprovalTask{
		Status:              "", // Invalid
		DocumentID:          1,
		ApplicationDocument: validDoc,
	}

	err := validators.ValidateStruct(&task)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}

// case: DocumentID is zero
func TestApprovalTaskValidation_DocumentIDZero(t *testing.T) {
	g := NewWithT(t)

	task := entity.ApprovalTask{
		Status:              "pending",
		DocumentID:          0, // Invalid
		ApplicationDocument: validDoc,
	}

	err := validators.ValidateStruct(&task)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}

// case: Document FileName is empty
func TestApprovalTaskValidation_DocFileNameEmpty(t *testing.T) {
	g := NewWithT(t)

	invalidDoc := validDoc
	invalidDoc.FileName = "" // Invalid

	task := entity.ApprovalTask{
		Status:              "pending",
		DocumentID:          1,
		ApplicationDocument: invalidDoc,
	}

	err := validators.ValidateStruct(&task)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}

// case: Document ApplicationScholarshipID is zero
func TestApprovalTaskValidation_DocAppScholarshipIDZero(t *testing.T) {
	g := NewWithT(t)

	invalidDoc := validDoc
	invalidDoc.ApplicationScholarshipID = 0 // Invalid

	task := entity.ApprovalTask{
		Status:              "pending",
		DocumentID:          1,
		ApplicationDocument: invalidDoc,
	}

	err := validators.ValidateStruct(&task)

	// Expect an error
	g.Expect(err).ToNot(BeNil())
}