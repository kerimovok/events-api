package controllers

import (
	"events-api/internal/requests"
	"events-api/internal/services"
	"events-api/pkg/utils"
	"events-api/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type EventController struct {
	eventService services.EventService
}

func NewEventController(eventService services.EventService) *EventController {
	return &EventController{
		eventService: eventService,
	}
}

func (ec *EventController) CreateEvent(c *fiber.Ctx) error {
	ctx := c.Context()
	var input requests.CreateEventRequest

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := validator.ValidateStruct(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	event, err := ec.eventService.CreateEvent(ctx, input.Name, input.Properties)
	if err != nil {
		utils.LogError("failed to create event", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create event", err)
	}

	return utils.SuccessResponse(c, "Event created successfully", event)
}
