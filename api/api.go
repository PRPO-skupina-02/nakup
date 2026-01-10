package api

import (
	"net/http"

	"github.com/PRPO-skupina-02/common/middleware"
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
//	@BasePath	/api/v1

func Register(router *gin.Engine, db *gorm.DB, trans ut.Translator) {
	// Healthcheck
	router.GET("/healthcheck", healthcheck)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// REST API
	v1 := router.Group("/api/v1")
	v1.Use(middleware.TransactionMiddleware(db))
	v1.Use(middleware.TranslationMiddleware(trans))
	v1.Use(middleware.ErrorMiddleware)

	// Reservations

}

func healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
