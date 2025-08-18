package main

import (
	"context"
	"events-api/internal/config"
	"events-api/internal/constants"
	"events-api/internal/database"
	"events-api/internal/queue"
	"events-api/internal/routes"
	"log"
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
	"github.com/kerimovok/go-pkg-database/mongo"
	pkgConfig "github.com/kerimovok/go-pkg-utils/config"
	pkgValidator "github.com/kerimovok/go-pkg-utils/validator"
)

func init() {
	// Load all configs
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("failed to load configs: %v", err)
	}

	// Validate environment variables
	if err := pkgValidator.ValidateConfig(constants.EnvValidationRules); err != nil {
		log.Fatalf("configuration validation failed: %v", err)
	}
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
	// Connect to MongoDB using go-pkg-database
	mongoConfig := mongo.MongoConfig{
		URI:            pkgConfig.GetEnv("DB_URI"),
		DBName:         pkgConfig.GetEnv("DB_NAME"),
		Timeout:        10 * time.Second,
		MaxPoolSize:    100,
		MinPoolSize:    5,
		MaxIdleTime:    5 * time.Minute,
		MaxConnecting:  10,
		ReadPreference: "primary",
		RetryWrites:    true,
		RetryReads:     true,
	}

	client, err := mongo.Connect(mongoConfig)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Set the global DB client for handlers
	database.DBClient = client

	// Get service configuration
	eventProcessingMode := pkgConfig.GetEnv("EVENT_PROCESSING_MODE")
	enableRestAPI := eventProcessingMode == "rest-only" || eventProcessingMode == "hybrid"
	enableRabbitMQConsumer := eventProcessingMode == "queue-only" || eventProcessingMode == "hybrid"

	log.Printf("Event processing mode: %s", eventProcessingMode)
	log.Printf("Service configuration: REST API=%v, RabbitMQ Consumer=%v", enableRestAPI, enableRabbitMQConsumer)

	// Validate processing mode
	if eventProcessingMode != "rest-only" && eventProcessingMode != "queue-only" && eventProcessingMode != "hybrid" {
		log.Fatal("Invalid EVENT_PROCESSING_MODE. Must be 'rest-only', 'queue-only', or 'hybrid'")
	}

	var app *fiber.App
	var consumer *queue.Consumer

	// Setup Fiber app only if REST API is enabled
	if enableRestAPI {
		app = setupApp()
		routes.Setup(app)
		log.Println("REST API server initialized")
	}

	// Setup RabbitMQ consumer only if enabled
	if enableRabbitMQConsumer {
		var err error
		consumer, err = queue.NewConsumer()
		if err != nil {
			log.Printf("Failed to initialize RabbitMQ consumer: %v", err)
			log.Println("Continuing without RabbitMQ consumer...")
			enableRabbitMQConsumer = false
		} else {
			// Start consuming messages in background
			go func() {
				if err := consumer.StartConsuming(); err != nil {
					log.Printf("RabbitMQ consumer error: %v", err)
				}
			}()
			log.Println("RabbitMQ consumer initialized")
		}
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Gracefully shutting down...")

		// Shutdown the server if REST API is enabled
		if enableRestAPI && app != nil {
			if err := app.Shutdown(); err != nil {
				log.Printf("error during server shutdown: %v", err)
			}
		}

		// Close RabbitMQ consumer if enabled
		if enableRabbitMQConsumer && consumer != nil {
			if err := consumer.Close(); err != nil {
				log.Printf("error during consumer shutdown: %v", err)
			}
		}

		log.Println("Server gracefully stopped")
		os.Exit(0)
	}()

	// Start server only if REST API is enabled
	if enableRestAPI {
		log.Printf("Starting REST API server on port %s", pkgConfig.GetEnv("PORT"))
		if err := app.Listen(":" + pkgConfig.GetEnv("PORT")); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	} else {
		log.Println("REST API is disabled, running in consumer-only mode")
		// Keep the main goroutine alive for the consumer
		select {}
	}
}
