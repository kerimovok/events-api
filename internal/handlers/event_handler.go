package handlers

import (
	"context"
	"events-api/internal/database"
	"events-api/internal/models"
	"events-api/internal/requests"
	internalUtils "events-api/internal/utils"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kerimovok/go-pkg-utils/httpx"
	"github.com/kerimovok/go-pkg-utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateEvent(c *fiber.Ctx) error {
	ctx := c.Context()
	var input requests.CreateEventRequest

	if err := c.BodyParser(&input); err != nil {
		response := httpx.BadRequest("Invalid request body", err)
		return httpx.SendResponse(c, response)
	}

	validationErrors := validator.ValidateStruct(&input)
	if validationErrors.HasErrors() {
		// Convert validator.ValidationErrors to []httpx.ValidationError
		httpxErrors := make([]httpx.ValidationError, len(validationErrors))
		for i, err := range validationErrors {
			httpxErrors[i] = httpx.ValidationError{
				Field:   err.Field,
				Message: err.Message,
			}
		}
		response := httpx.UnprocessableEntityWithValidation("Validation failed", httpxErrors)
		return httpx.SendValidationResponse(c, response)
	}

	event := models.Event{
		Id:         primitive.NewObjectID(),
		Properties: input.Properties,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result, err := database.DBClient.Database().Collection("events").InsertOne(ctx, event)
	if err != nil {
		log.Printf("failed to create event: %v", err)
		response := httpx.InternalServerError("Internal server error", err)
		return httpx.SendResponse(c, response)
	}

	event.Id = result.InsertedID.(primitive.ObjectID)

	response := httpx.Created("Event created successfully", event)
	return httpx.SendResponse(c, response)
}

// GetEvents retrieves a paginated list of events with optional filtering and sorting
// Supports query parameters:
// - page: Page number (default: 1)
// - limit: Items per page (default: 10)
// - sortBy: Field to sort by (default: createdAt)
// - sortOrder: Sort direction, 'asc' or 'desc' (default: asc)
// - Any other query parameter will be used as a filter
func GetEvents(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Extract query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "createdAt")
	sortOrder := c.Query("sortOrder", "asc")

	// Extract filters
	filters := bson.M{}
	for key, values := range c.Queries() {
		if key != "page" && key != "limit" && key != "sortBy" && key != "sortOrder" {
			filters[key] = values
		}
	}

	// Query events
	events, err := internalUtils.QueryEvents(ctx, filters, internalUtils.BuildSortOptions(sortBy, sortOrder), page, limit)
	if err != nil {
		return httpx.SendResponse(c, httpx.InternalServerError("Failed to fetch events", err))
	}

	return httpx.SendResponse(c, httpx.OK("Events retrieved successfully", fiber.Map{
		"page":   page,
		"limit":  limit,
		"events": events,
	}))
}

// GetStats aggregates event data based on grouping and aggregation criteria
// Supports query parameters:
// - groupBy: Field to group results by
// - aggregates: Aggregation operation (count, sum, avg)
// - Any other query parameter will be used as a filter
func GetStats(c *fiber.Ctx) error {
	ctx := context.Background()

	// Extract query parameters
	groupBy := c.Query("groupBy", "")
	aggregates := c.Query("aggregates", "count") // E.g., count, sum, avg
	filters := bson.M{}
	for key, values := range c.Queries() {
		if key != "groupBy" && key != "aggregates" {
			filters[key] = values
		}
	}

	// Perform aggregation query
	stats, err := internalUtils.AggregateStats(ctx, filters, groupBy, aggregates)
	if err != nil {
		return httpx.SendResponse(c, httpx.InternalServerError("Failed to fetch stats", err))
	}

	return httpx.SendResponse(c, httpx.OK("Stats retrieved successfully", fiber.Map{
		"groupBy":    groupBy,
		"aggregates": aggregates,
		"stats":      stats,
	}))
}

// GetTimeSeries generates time-based aggregations of event data
// Supports query parameters:
// - interval: Time grouping interval (hour, day, week, month)
// - aggregates: Aggregation operation (count, sum, avg)
// - Any other query parameter will be used as a filter
func GetTimeSeries(c *fiber.Ctx) error {
	ctx := context.Background()

	// Extract query parameters
	aggregates := c.Query("aggregates", "count")
	interval := c.Query("interval", "day") // E.g., day, week, month
	filters := bson.M{}
	for key, values := range c.Queries() {
		if key != "aggregates" && key != "interval" {
			filters[key] = values
		}
	}

	// Perform time-series query
	timeSeries, err := internalUtils.AggregateTimeSeries(ctx, filters, interval, aggregates)
	if err != nil {
		return httpx.SendResponse(c, httpx.InternalServerError("Failed to fetch time series", err))
	}

	return httpx.SendResponse(c, httpx.OK("Time series retrieved successfully", fiber.Map{
		"interval":   interval,
		"aggregates": aggregates,
		"timeSeries": timeSeries,
	}))
}
