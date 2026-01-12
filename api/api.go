package api

import (
	"net/http"

	"github.com/PRPO-skupina-02/common/middleware"
	_ "github.com/PRPO-skupina-02/nakup/api/docs"
	"github.com/PRPO-skupina-02/nakup/services"
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
	reservations := v1.Group("/reservations/:reservationID")
	reservations.Use(ReservationContextMiddleware)

	v1.GET("/reservations", ReservationsList)
	reservations.GET("", ReservationsShow)
	v1.POST("/reservations", ReservationsCreate)
	reservations.PUT("", ReservationsUpdate)
	reservations.DELETE("", ReservationsDelete)

	// Purchases
	reservations.GET("/purchases", PurchasesList)
	reservations.GET("/purchases/:purchaseID", PurchasesShow)
	reservations.POST("/purchases", PurchasesCreate)
	reservations.PUT("/purchases/:purchaseID", PurchasesUpdate)
	reservations.DELETE("/purchases/:purchaseID", PurchasesDelete)
}

func healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
