package routes

import (
	"events-api/internal/handlers"

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
	event := v1.Group("/events")
	event.Post("/", handlers.CreateEvent)
	event.Get("/", handlers.GetEvents)
	event.Get("/stats", handlers.GetStats)
	event.Get("/timeseries", handlers.GetTimeSeries)
}
