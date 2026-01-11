package api

import (
	"net/http"
	"time"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/request"
	"github.com/PRPO-skupina-02/nakup/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReservationResponse struct {
	ID         uuid.UUID              `json:"id"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	TimeSlotID uuid.UUID              `json:"time_slot_id"`
	UserID     uuid.UUID              `json:"user_id"`
	Type       models.ReservationType `json:"type"`
	Row        int                    `json:"row"`
	Col        int                    `json:"col"`
}

func newReservationResponse(reservation models.Reservation) ReservationResponse {
	return ReservationResponse{
		ID:         reservation.ID,
		CreatedAt:  reservation.CreatedAt,
		UpdatedAt:  reservation.UpdatedAt,
		TimeSlotID: reservation.TimeSlotID,
		UserID:     reservation.UserID,
		Type:       reservation.Type,
		Row:        reservation.Row,
		Col:        reservation.Col,
	}
}

// ReservationsList
//
//	@Id				ReservationsList
//	@Summary		List reservations
//	@Description	List reservations
//	@Tags			reservations
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit the number of responses"	Default(10)
//	@Param			offset	query		int		false	"Offset the first response"		Default(0)
//	@Param			sort	query		string	false	"Sort results"
//	@Success		200		{object}	[]ReservationResponse
//	@Failure		400		{object}	middleware.HttpError
//	@Failure		404		{object}	middleware.HttpError
//	@Failure		500		{object}	middleware.HttpError
//	@Router			/reservations [get]
func ReservationsList(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	pagination := request.GetNormalizedPaginationArgs(c)
	sort := request.GetSortOptions(c)

	reservations, total, err := models.GetReservations(tx, pagination, sort)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := []ReservationResponse{}

	for _, reservation := range reservations {
		response = append(response, newReservationResponse(reservation))
	}

	request.RenderPaginatedResponse(c, response, total)
}

type ReservationRequest struct {
	TimeSlotID uuid.UUID              `json:"time_slot_id" binding:"required"`
	TheaterID  uuid.UUID              `json:"theater_id" binding:"required"`
	RoomID     uuid.UUID              `json:"room_id" binding:"required"`
	Type       models.ReservationType `json:"type" binding:"required,oneof=ONLINE POS"`
	Row        int                    `json:"row" binding:"required,min=1"`
	Col        int                    `json:"col" binding:"required,min=1"`
}

// ReservationsCreate
//
//	@Id				ReservationsCreate
//	@Summary		Create reservation
//	@Description	Create reservation
//	@Tags			reservations
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ReservationRequest	true	"request body"
//	@Success		200		{object}	ReservationResponse
//	@Failure		400		{object}	middleware.HttpError
//	@Failure		404		{object}	middleware.HttpError
//	@Failure		500		{object}	middleware.HttpError
//	@Router			/reservations [post]
func ReservationsCreate(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	validator := GetTimeSlotValidator(c)

	var req ReservationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = validator.ValidateTimeSlotExists(req.TheaterID, req.RoomID, req.TimeSlotID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	reservation := models.Reservation{
		ID:         uuid.New(),
		TimeSlotID: req.TimeSlotID,
		UserID:     uuid.Nil,
		Type:       req.Type,
		Row:        req.Row,
		Col:        req.Col,
	}

	err = reservation.Create(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, newReservationResponse(reservation))
}

// ReservationsShow
//
//	@Id				ReservationsShow
//	@Summary		Show reservation
//	@Description	Show reservation
//	@Tags			reservations
//	@Accept			json
//	@Produce		json
//	@Param			reservationID	path		string	true	"Reservation ID"	Format(uuid)
//	@Success		200				{object}	ReservationResponse
//	@Failure		400				{object}	middleware.HttpError
//	@Failure		404				{object}	middleware.HttpError
//	@Failure		500				{object}	middleware.HttpError
//	@Router			/reservations/{reservationID} [get]
func ReservationsShow(c *gin.Context) {
	reservation := GetContextReservation(c)
	c.JSON(http.StatusOK, newReservationResponse(reservation))
}

// ReservationsUpdate
//
//	@Id				ReservationsUpdate
//	@Summary		Update reservation
//	@Description	Update reservation
//	@Tags			reservations
//	@Accept			json
//	@Produce		json
//	@Param			reservationID	path		string				true	"Reservation ID"	Format(uuid)
//	@Param			request			body		ReservationRequest	true	"request body"
//	@Success		200				{object}	ReservationResponse
//	@Failure		400				{object}	middleware.HttpError
//	@Failure		404				{object}	middleware.HttpError
//	@Failure		500				{object}	middleware.HttpError
//	@Router			/reservations/{reservationID} [put]
func ReservationsUpdate(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	validator := GetTimeSlotValidator(c)
	reservation := GetContextReservation(c)

	var req ReservationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = validator.ValidateTimeSlotExists(req.TheaterID, req.RoomID, req.TimeSlotID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	reservation.TimeSlotID = req.TimeSlotID
	reservation.Type = req.Type
	reservation.Row = req.Row
	reservation.Col = req.Col

	err = reservation.Save(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, newReservationResponse(reservation))
}

// ReservationsDelete
//
//	@Id				ReservationsDelete
//	@Summary		Delete reservation
//	@Description	Delete reservation
//	@Tags			reservations
//	@Accept			json
//	@Produce		json
//	@Param			reservationID	path	string	true	"Reservation ID"	Format(uuid)
//	@Success		204
//	@Failure		400	{object}	middleware.HttpError
//	@Failure		404	{object}	middleware.HttpError
//	@Failure		500	{object}	middleware.HttpError
//	@Router			/reservations/{reservationID} [delete]
func ReservationsDelete(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	reservation := GetContextReservation(c)

	err := models.DeleteReservation(tx, reservation.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, "")
}
