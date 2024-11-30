package controllers

import (
	"events-api/internal/services"
	"events-api/pkg/errors"
	"events-api/pkg/request"
	"events-api/pkg/utils"

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
	var input request.CreateEventRequest

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, errors.NewBadRequestError("Invalid request body", err))
	}

	if err := utils.ValidateRequest(input); err != nil {
		return utils.ErrorResponse(c, errors.NewBadRequestError("Validation failed", err))
	}

	event, err := ec.eventService.CreateEvent(ctx, input.Name, input.Properties)
	if err != nil {
		return utils.ErrorResponse(c, errors.NewInternalError("Failed to create event", err))
	}

	return utils.SuccessResponse(c, "Event created successfully", event)
}
