package services

import (
	"errors"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/google/uuid"
)

type MockTimeSlotValidator struct {
	ValidTimeSlots map[string]bool
	ShouldError    bool
	Error          error
}

func NewMockTimeSlotValidator() *MockTimeSlotValidator {
	return &MockTimeSlotValidator{
		ValidTimeSlots: make(map[string]bool),
		ShouldError:    false,
	}
}

func (m *MockTimeSlotValidator) AddValidTimeSlot(theaterID, roomID, timeSlotID uuid.UUID) {
	key := m.makeKey(theaterID, roomID, timeSlotID)
	m.ValidTimeSlots[key] = true
}

func (m *MockTimeSlotValidator) makeKey(theaterID, roomID, timeSlotID uuid.UUID) string {
	return theaterID.String() + "|" + roomID.String() + "|" + timeSlotID.String()
}

func (m *MockTimeSlotValidator) ValidateTimeSlotExists(theaterID, roomID, timeSlotID uuid.UUID) error {
	if m.ShouldError {
		if m.Error != nil {
			return m.Error
		}
		return errors.New("mock error")
	}

	key := m.makeKey(theaterID, roomID, timeSlotID)
	if !m.ValidTimeSlots[key] {
		return middleware.NewNotFoundError()
	}

	return nil
}
