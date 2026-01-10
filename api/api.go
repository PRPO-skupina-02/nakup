package api

import (
	"net/http"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/nakup/services"
	_ "github.com/PRPO-skupina-02/spored/api/docs"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

//	@title			Nakup API
//	@version		1.0
//	@description	API za upravljanje z kinodvoranami in njihovim sporedom

//	@host		localhost:8081
//	@BasePath	/api/v1/nakup

func Register(router *gin.Engine, db *gorm.DB, trans ut.Translator, timeSlotValidator services.TimeSlotValidator) {
	// Healthcheck
	router.GET("/healthcheck", healthcheck)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// REST API
	v1 := router.Group("/api/v1/nakup")
	v1.Use(middleware.TransactionMiddleware(db))
	v1.Use(middleware.TranslationMiddleware(trans))
	v1.Use(middleware.ErrorMiddleware)
	v1.Use(TimeSlotValidatorMiddleware(timeSlotValidator))

	// Reservations

	v1.GET("/reservations", ReservationsList)
	v1.POST("/reservations", ReservationsCreate)
	v1.GET("/reservations/:reservationID", ReservationsShow)
	v1.PUT("/reservations/:reservationID", ReservationsUpdate)
	v1.DELETE("/reservations/:reservationID", ReservationsDelete)
}

func healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
