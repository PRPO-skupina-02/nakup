package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/PRPO-skupina-02/common/config"
	"github.com/PRPO-skupina-02/common/database"
	"github.com/PRPO-skupina-02/common/logging"
	"github.com/PRPO-skupina-02/common/validation"
	"github.com/PRPO-skupina-02/nakup/api"
	"github.com/PRPO-skupina-02/nakup/clients/spored/client"
	"github.com/PRPO-skupina-02/nakup/db"
	"github.com/PRPO-skupina-02/nakup/services"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {
	slog.Info("Starting server")

	logger := logging.GetDefaultLogger()
	slog.SetDefault(logger)

	db, err := database.OpenAndMigrateProd(db.MigrationsFS)
	if err != nil {
		return err
	}

	trans, err := validation.RegisterValidation()
	if err != nil {
		return err
	}

	sporedHost := config.GetEnv("SPORED_HOST")
	transportConfig := client.DefaultTransportConfig().WithHost(sporedHost)
	sporedClient := client.NewHTTPClientWithConfig(strfmt.Default, transportConfig)

	timeSlotService := services.NewSporedTimeSlotService(sporedClient)

	authHost := config.GetEnv("AUTH_HOST")

	router := gin.Default()
	api.Register(router, db, trans, timeSlotService, authHost)

	slog.Info("Server startup complete")
	err = router.Run(":8081")
	if err != nil {
		return err
	}

	return nil
}
