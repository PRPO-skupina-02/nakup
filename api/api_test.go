package api

import (
	"testing"

	"github.com/PRPO-skupina-02/common/clients/auth/models"
	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/validation"
	"github.com/PRPO-skupina-02/nakup/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// MockUserMiddleware creates a test middleware that sets a mock user in the context
func MockUserMiddleware(userID uuid.UUID, role models.ModelsUserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &models.APIUserResponse{
			ID:        userID.String(),
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Role:      role,
			Active:    true,
		}
		middleware.SetContextUser(c, user)
		c.Next()
	}
}

func TestingRouter(t *testing.T, db *gorm.DB, timeSlotService services.TimeSlotService) *gin.Engine {
	router := gin.Default()
	trans, err := validation.RegisterValidation()
	require.NoError(t, err)

	// Use mock auth instead of real auth for testing
	router.Use(MockUserMiddleware(uuid.MustParse("00000000-0000-0000-0000-000000000001"), models.ModelsUserRoleEmployee))

	// Register routes but skip the auth middleware since we added mock above
	v1 := router.Group("/api/v1/nakup")
	v1.Use(middleware.TransactionMiddleware(db))
	v1.Use(middleware.TranslationMiddleware(trans))
	v1.Use(middleware.ErrorMiddleware)
	v1.Use(TimeSlotServiceMiddleware(timeSlotService))

	// Reservations
	v1.POST("/reservations", ReservationsCreate)
	v1.GET("/reservations/my", MyReservationsList)
	v1.GET("/reservations", ReservationsList)

	reservations := v1.Group("/reservations/:reservationID")
	reservations.Use(ReservationContextMiddleware)
	reservations.GET("", ReservationsShow)
	reservations.PUT("", ReservationsUpdate)
	reservations.DELETE("", ReservationsDelete)

	// Purchases
	purchases := v1.Group("/reservations/:reservationID/purchases")
	purchases.Use(ReservationContextMiddleware)
	purchases.GET("", PurchasesList)
	purchases.GET("/:purchaseID", PurchasesShow)
	purchases.POST("", PurchasesCreate)
	purchases.PUT("/:purchaseID", PurchasesUpdate)
	purchases.DELETE("/:purchaseID", PurchasesDelete)

	return router
}
