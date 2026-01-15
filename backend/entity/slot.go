package entity
import (
	"time"
	"gorm.io/gorm"
)

type Slot struct {
	gorm.Model
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	Capacity uint   `gorm:"not null" json:"capacity"`
	BookCount uint   `gorm:"not null" json:"book_count"`
	Status 	string `gorm:"not null" json:"status"`
	IsBooked bool `gorm:"not null;default:false" json:"is_booked"`

	InterviewRoundID uint `json:"interview_round_id"`
	InterviewRound   InterviewRound `gorm:"foreignKey:InterviewRoundID" json:"-"`

		InterviewerSlots  []InterviewerSlot  `gorm:"foreignKey:SlotID" json:"interviewer_slots"`

		IntervieweBookings []IntervieweBooking `gorm:"foreignKey:SlotID" json:"interviewe_bookings"`

	}

	

	// AfterFind GORM hook to ensure slices are never nil, so they marshal as [] instead of null

	func (slot *Slot) AfterFind(tx *gorm.DB) (err error) {

		if slot.InterviewerSlots == nil {

			slot.InterviewerSlots = []InterviewerSlot{}

		}

		if slot.IntervieweBookings == nil {

			slot.IntervieweBookings = []IntervieweBooking{}

		}

		return

	}
