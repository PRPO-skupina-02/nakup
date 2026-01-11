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
	TimeSlotValidatorKey  = "timeslot_validator"
	contextReservationKey = "reservation"
)

func TimeSlotValidatorMiddleware(validator services.TimeSlotValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(TimeSlotValidatorKey, validator)
		c.Next()
	}
}

func GetTimeSlotValidator(c *gin.Context) services.TimeSlotValidator {
	timeSlotValidator, exists := c.Get(TimeSlotValidatorKey)
	if !exists {
		return nil
	}
	return timeSlotValidator.(services.TimeSlotValidator)
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
