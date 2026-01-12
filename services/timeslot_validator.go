package services

import (
	"log/slog"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/nakup/clients/spored/client/timeslots"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type TimeSlotValidator interface {
	ValidateTimeSlotExists(theaterID, roomID, timeSlotID uuid.UUID) error
}

type SporedTimeSlotValidator struct {
	client timeslots.ClientService
}

func NewSporedTimeSlotValidator(client timeslots.ClientService) TimeSlotValidator {
	return &SporedTimeSlotValidator{
		client: client,
	}
}

func (v *SporedTimeSlotValidator) ValidateTimeSlotExists(theaterID, roomID, timeSlotID uuid.UUID) error {
	params := timeslots.NewTimeSlotsShowParams()
	params.TheaterID = strfmt.UUID(theaterID.String())
	params.RoomID = strfmt.UUID(roomID.String())
	params.TimeSlotID = strfmt.UUID(timeSlotID.String())

	_, err := v.client.TimeSlotsShow(params)
	if err != nil {
		slog.Error("failed to fetch timeslot", "err", err)
		return middleware.NewNotFoundError()
	}

	return nil
}
