package routes

import (
	"events-api/internal/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Setup(app *fiber.App) {
	// API routes group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Monitor route
	app.Get("/metrics", monitor.New())

	// Event routes
	eventRoutes := v1.Group("/events")
	eventRoutes.Post("/", controllers.CreateEvent)
}
