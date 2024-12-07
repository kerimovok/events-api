package handlers

import (
	"events-api/internal/models"
	"events-api/internal/requests"
	"events-api/pkg/database"
	"events-api/pkg/utils"
	"events-api/pkg/validator"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateEvent(c *fiber.Ctx) error {
	ctx := c.Context()
	var input requests.CreateEventRequest

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := validator.ValidateStruct(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	event := models.Event{
		Id:         primitive.NewObjectID(),
		Properties: input.Properties,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result, err := database.DBClient.Database(utils.GetEnv("DB_NAME")).Collection("events").InsertOne(ctx, event)
	if err != nil {
		utils.LogError("failed to create event", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error", err)
	}

	event.Id = result.InsertedID.(primitive.ObjectID)

	return utils.SuccessResponse(c, "Event created successfully", event)
}
