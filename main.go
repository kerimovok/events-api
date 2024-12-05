package main

import (
	"context"
	"events-api/internal/routes"
	"events-api/pkg/config"
	"events-api/pkg/database"
	"events-api/pkg/validator"
	"log"

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
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}

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
	app := setupApp()

	dbService := database.NewMongoService()
	if err := dbService.Connect(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := dbService.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	routes.Setup(app, dbService)

	log.Fatal(app.Listen(":" + config.Env.Server.Port))
}
