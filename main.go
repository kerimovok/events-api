package main

import (
	"context"
	"events-api/internal/config"
	"events-api/internal/constants"
	"events-api/internal/routes"
	"events-api/pkg/database"
	"events-api/pkg/utils"
	"events-api/pkg/validator"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

func init() {
	// Load all configs
	if err := config.LoadConfig(); err != nil {
		utils.LogFatal("failed to load configs", err)
	}

	// Validate environment variables
	if err := utils.ValidateConfig(constants.EnvValidationRules); err != nil {
		utils.LogFatal("configuration validation failed", err)
	}

	// Initialize validator
	validator.InitValidator()
}

func setupApp() *fiber.App {
	app := fiber.New(fiber.Config{})

	// Middleware
	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(healthcheck.New())
	app.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return uuid.New().String()
		},
	}))
	app.Use(logger.New())

	return app
}

func main() {
	// Connect to MongoDB
	if err := database.ConnectDB(); err != nil {
		utils.LogFatal("failed to connect to MongoDB", err)
	}

	app := setupApp()
	routes.Setup(app)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		utils.LogInfo("Gracefully shutting down...")

		// Shutdown the server
		if err := app.Shutdown(); err != nil {
			utils.LogError("error during server shutdown", err)
		}

		// Disconnect from MongoDB
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := database.DisconnectDB(ctx); err != nil {
			utils.LogError("failed to disconnect from MongoDB", err)
		}

		utils.LogInfo("Server gracefully stopped")
		os.Exit(0)
	}()

	// Start server
	if err := app.Listen(":" + utils.GetEnv("PORT")); err != nil && err != http.ErrServerClosed {
		utils.LogFatal("failed to start server", err)
	}
}
