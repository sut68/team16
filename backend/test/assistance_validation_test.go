package test

import (
	"backend/entity"
	"backend/validators"
	"testing"

	. "github.com/onsi/gomega"
)

func TestAssistanceFullValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	Assistance := entity.Assistance{
		Massage:    "hello",
		ChatroomID: 1,
		SenderID:   1,
	}

	err := validators.ValidateStruct(&Assistance)

	g.Expect(err).To(BeNil())
}

func TestMassageValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("massage is required", func(t *testing.T) {
		s := entity.Assistance{
			Massage:    "",
			ChatroomID: 1,
			SenderID:   1,
		}

		err := validators.ValidateStruct(s)

		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Error()).To(Equal("massage is required"))
	})
}

func TestChatroomValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("chatroom is required", func(t *testing.T) {
		s := entity.Assistance{
			Massage:    "hello",
			ChatroomID: 0,
			SenderID:   1,
		}

		err := validators.ValidateStruct(s)

		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Error()).To(Equal("chatroom is required"))
	})
}

func TestSenderValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("sender is required", func(t *testing.T) {
		s := entity.Assistance{
			Massage:    "hello",
			ChatroomID: 1,
			SenderID:   0,
		}

		err := validators.ValidateStruct(s)

		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Error()).To(Equal("sender is required"))
	})
}

