package api

import (
	"errors"
	"net/http"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/request"
	"github.com/PRPO-skupina-02/nakup/models"
	"github.com/PRPO-skupina-02/nakup/services"
	"github.com/gin-gonic/gin"
)

const (
	TimeSlotServiceKey    = "timeslot_service"
	contextReservationKey = "reservation"
)

func TimeSlotServiceMiddleware(service services.TimeSlotService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(TimeSlotServiceKey, service)
		c.Next()
	}
}

func GetTimeSlotService(c *gin.Context) services.TimeSlotService {
	timeSlotService, exists := c.Get(TimeSlotServiceKey)
	if !exists {
		return nil
	}
	return timeSlotService.(services.TimeSlotService)
}

func SetContextReservation(c *gin.Context, reservation models.Reservation) {
	c.Set(contextReservationKey, reservation)
}

func GetContextReservation(c *gin.Context) models.Reservation {
	reservation, ok := c.Get(contextReservationKey)
	if !ok {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("Could not get reservation from context"))
		return models.Reservation{}
	}

	return reservation.(models.Reservation)
}

func ReservationContextMiddleware(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	id, err := request.GetUUIDParam(c, "reservationID")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	reservation, err := models.GetReservation(tx, id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	SetContextReservation(c, reservation)

	c.Next()
}
