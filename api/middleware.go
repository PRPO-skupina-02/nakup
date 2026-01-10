package api

import (
	"github.com/PRPO-skupina-02/nakup/services"
	"github.com/gin-gonic/gin"
)

const TimeSlotValidatorKey = "timeslot_validator"

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
