package routes

import (
	"events-api/internal/controllers"
	"events-api/internal/services"
	"events-api/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Setup(app *fiber.App, dbService database.DatabaseService) {
	api := app.Group("/api/v1")

	// Services
	eventService := services.NewEventService(dbService)

	// Controllers
	eventController := controllers.NewEventController(eventService)

	// Monitor route
	app.Get("/metrics", monitor.New())

	// Event routes
	eventRoutes := api.Group("/events")
	eventRoutes.Post("/", eventController.CreateEvent)
}
