package api

import (
	"net/http"

	"github.com/PRPO-skupina-02/common/clients/auth/models"
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

func Register(router *gin.Engine, db *gorm.DB, trans ut.Translator, timeSlotService services.TimeSlotService, authHost string) {
	// Healthcheck
	router.GET("/healthcheck", healthcheck)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// REST API
	v1 := router.Group("/api/v1/nakup")
	v1.Use(middleware.TransactionMiddleware(db))
	v1.Use(middleware.TranslationMiddleware(trans))
	v1.Use(middleware.ErrorMiddleware)
	v1.Use(TimeSlotServiceMiddleware(timeSlotService))
	v1.Use(middleware.UserMiddleware(authHost))

	// Reservations
	v1.POST("/reservations", ReservationsCreate)
	v1.GET("/reservations/my", MyReservationsList)

	reservationsStaff := v1.Group("/reservations")
	reservationsStaff.Use(middleware.RequireRole(models.ModelsUserRoleEmployee, models.ModelsUserRoleAdmin))
	reservationsStaff.GET("", ReservationsList)

	reservations := v1.Group("/reservations/:reservationID")
	reservations.Use(ReservationContextMiddleware)
	reservations.Use(middleware.RequireRole(models.ModelsUserRoleEmployee, models.ModelsUserRoleAdmin))
	reservations.GET("", ReservationsShow)
	reservations.PUT("", ReservationsUpdate)
	reservations.DELETE("", ReservationsDelete)

	// Purchases
	purchases := v1.Group("/reservations/:reservationID/purchases")
	purchases.Use(ReservationContextMiddleware)
	purchases.Use(middleware.RequireRole(models.ModelsUserRoleEmployee, models.ModelsUserRoleAdmin))
	purchases.GET("", PurchasesList)
	purchases.GET("/:purchaseID", PurchasesShow)
	purchases.POST("", PurchasesCreate)
	purchases.PUT("/:purchaseID", PurchasesUpdate)
	purchases.DELETE("/:purchaseID", PurchasesDelete)
}

func healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
