package services

import (
	"log/slog"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/nakup/clients/spored/client"
	"github.com/PRPO-skupina-02/nakup/clients/spored/client/rooms"
	"github.com/PRPO-skupina-02/nakup/clients/spored/client/timeslots"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type TimeSlotInfo struct {
	TimeSlotID uuid.UUID
	RoomID     uuid.UUID
	TheaterID  uuid.UUID
	Rows       int
	Columns    int
}

type TimeSlotService interface {
	ValidateTimeSlotExists(theaterID, roomID, timeSlotID uuid.UUID) (*TimeSlotInfo, error)
}

type SporedTimeSlotService struct {
	timeslotClient timeslots.ClientService
	roomClient     rooms.ClientService
}

func NewSporedTimeSlotService(client *client.Spored) TimeSlotService {
	return &SporedTimeSlotService{
		timeslotClient: client.Timeslots,
		roomClient:     client.Rooms,
	}
}

func (v *SporedTimeSlotService) ValidateTimeSlotExists(theaterID, roomID, timeSlotID uuid.UUID) (*TimeSlotInfo, error) {
	params := timeslots.NewTimeSlotsShowParams()
	params.TheaterID = strfmt.UUID(theaterID.String())
	params.RoomID = strfmt.UUID(roomID.String())
	params.TimeSlotID = strfmt.UUID(timeSlotID.String())

	_, err := v.timeslotClient.TimeSlotsShow(params)
	if err != nil {
		slog.Error("failed to fetch timeslot", "err", err)
		return nil, middleware.NewNotFoundError()
	}

	roomParams := rooms.NewRoomsShowParams()
	roomParams.TheaterID = strfmt.UUID(theaterID.String())
	roomParams.RoomID = strfmt.UUID(roomID.String())

	roomResp, err := v.roomClient.RoomsShow(roomParams)
	if err != nil {
		return nil, middleware.NewNotFoundError()
	}

	return &TimeSlotInfo{
		TimeSlotID: timeSlotID,
		RoomID:     roomID,
		TheaterID:  theaterID,
		Rows:       int(roomResp.Payload.Rows),
		Columns:    int(roomResp.Payload.Columns),
	}, nil
}
