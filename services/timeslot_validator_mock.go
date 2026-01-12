package services

import (
	"errors"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/google/uuid"
)

type MockTimeSlotInfo struct {
	Rows    int
	Columns int
}

type MockTimeSlotService struct {
	ValidTimeSlots map[string]MockTimeSlotInfo
	ShouldError    bool
	Error          error
}

func NewMockTimeSlotService() *MockTimeSlotService {
	return &MockTimeSlotService{
		ValidTimeSlots: make(map[string]MockTimeSlotInfo),
		ShouldError:    false,
	}
}

func (m *MockTimeSlotService) AddValidTimeSlot(theaterID, roomID, timeSlotID uuid.UUID) {
	m.AddValidTimeSlotWithRoom(theaterID, roomID, timeSlotID, 10, 10)
}

func (m *MockTimeSlotService) AddValidTimeSlotWithRoom(theaterID, roomID, timeSlotID uuid.UUID, rows, columns int) {
	key := m.makeKey(theaterID, roomID, timeSlotID)
	m.ValidTimeSlots[key] = MockTimeSlotInfo{
		Rows:    rows,
		Columns: columns,
	}
}

func (m *MockTimeSlotService) makeKey(theaterID, roomID, timeSlotID uuid.UUID) string {
	return theaterID.String() + "|" + roomID.String() + "|" + timeSlotID.String()
}

func (m *MockTimeSlotService) ValidateTimeSlotExists(theaterID, roomID, timeSlotID uuid.UUID) (*TimeSlotInfo, error) {
	if m.ShouldError {
		if m.Error != nil {
			return nil, m.Error
		}
		return nil, errors.New("mock error")
	}

	key := m.makeKey(theaterID, roomID, timeSlotID)
	info, exists := m.ValidTimeSlots[key]
	if !exists {
		return nil, middleware.NewNotFoundError()
	}

	return &TimeSlotInfo{
		TimeSlotID: timeSlotID,
		RoomID:     roomID,
		TheaterID:  theaterID,
		Rows:       info.Rows,
		Columns:    info.Columns,
	}, nil
}
