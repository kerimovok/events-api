package handlers

import (
	"context"
	"encoding/json"
	"events-api/internal/constants"
	"events-api/internal/database"
	"events-api/internal/models"
	"events-api/internal/requests"
	internalUtils "events-api/internal/utils"
	"fmt"
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
		log.Printf("failed to parse request body: %v", err)
		response := httpx.BadRequest("Invalid request body", err)
		return httpx.SendResponse(c, response)
	}

	validationErrors := validator.ValidateStruct(&input)
	if validationErrors.HasErrors() {
		log.Printf("validation failed for event creation: %v", validationErrors)
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

	result, err := database.DBClient.Database().Collection(constants.EventsCollection).InsertOne(ctx, event)
	if err != nil {
		log.Printf("failed to create event in database: %v", err)
		response := httpx.InternalServerError("Failed to create event", err)
		return httpx.SendResponse(c, response)
	}

	event.Id = result.InsertedID.(primitive.ObjectID)
	log.Printf("event created successfully with ID: %s", event.Id.Hex())

	response := httpx.Created("Event created successfully", event)
	return httpx.SendResponse(c, response)
}

// GetEvents retrieves a paginated list of events with optional filtering and sorting
// Supports query parameters:
// - page: Page number (default: 1)
// - limit: Items per page (default: 50)
// - sortBy: Field to sort by (default: createdAt)
// - sortOrder: Sort direction, 'asc' or 'desc' (default: asc)
// - filters: JSON string for complex MongoDB queries (e.g., {"status":"active","created_at":{"$gte":"2024-01-01"}})
// - Any other query parameter will be used as a simple filter (for backward compatibility)
func GetEvents(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.QueryTimeout)
	defer cancel()

	// Extract query parameters
	page, err := strconv.Atoi(c.Query(constants.ParamPage, strconv.Itoa(constants.DefaultPage)))
	if err != nil {
		return httpx.SendResponse(c, httpx.BadRequest("Invalid page parameter", err))
	}
	limit, err := strconv.Atoi(c.Query(constants.ParamLimit, strconv.Itoa(constants.DefaultLimit)))
	if err != nil {
		return httpx.SendResponse(c, httpx.BadRequest("Invalid limit parameter", err))
	}
	sortBy := c.Query(constants.ParamSortBy, constants.DefaultSortBy)
	sortOrder := c.Query(constants.ParamSortOrder, constants.DefaultSortOrder)

	// Validate parameters
	if page < 1 {
		return httpx.SendResponse(c, httpx.BadRequest("Page must be a positive number", nil))
	}
	if limit < 1 || limit > 1000 {
		return httpx.SendResponse(c, httpx.BadRequest("Limit must be between 1 and 1000", nil))
	}
	if !isValidSortField(sortBy) {
		return httpx.SendResponse(c, httpx.BadRequest("Invalid sort field", nil))
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		return httpx.SendResponse(c, httpx.BadRequest("Sort order must be 'asc' or 'desc'", nil))
	}

	// Parse JSON filters if provided
	var filters bson.M
	if filterStr := c.Query("filters"); filterStr != "" {
		var err error
		filters, err = parseJSONFilters(filterStr)
		if err != nil {
			return httpx.SendResponse(c, httpx.BadRequest("Invalid filters parameter", err))
		}
	} else {
		// Extract simple filters for backward compatibility
		filters = bson.M{}
		for key, values := range c.Queries() {
			if !isReservedQueryParam(key) {
				filters[key] = values
			}
		}
	}

	// Query events
	events, err := internalUtils.QueryEvents(ctx, filters, internalUtils.BuildSortOptions(sortBy, sortOrder), page, limit)
	if err != nil {
		return httpx.SendResponse(c, httpx.InternalServerError("Failed to fetch events", err))
	}

	return httpx.SendResponse(c, httpx.OK("Events retrieved successfully", fiber.Map{
		constants.ParamPage:  page,
		constants.ParamLimit: limit,
		"events":             events,
	}))
}

// isReservedQueryParam checks if a query parameter is reserved for pagination/sorting
func isReservedQueryParam(key string) bool {
	reservedParams := []string{constants.ParamPage, constants.ParamLimit, constants.ParamSortBy, constants.ParamSortOrder, constants.ParamFilters}
	for _, reserved := range reservedParams {
		if key == reserved {
			return true
		}
	}
	return false
}

// isValidSortField checks if a sort field is valid to prevent injection attacks
func isValidSortField(field string) bool {
	validFields := []string{"created_at", "updated_at", "id"}
	for _, valid := range validFields {
		if field == valid {
			return true
		}
	}
	return false
}

// isValidAggregation checks if an aggregation type is valid
func isValidAggregation(agg string) bool {
	validAggregations := []string{"count", "sum", "avg"}
	for _, valid := range validAggregations {
		if agg == valid {
			return true
		}
	}
	return false
}

// isValidTimeInterval checks if a time interval is valid
func isValidTimeInterval(interval string) bool {
	validIntervals := []string{"hour", "day", "week", "month"}
	for _, valid := range validIntervals {
		if interval == valid {
			return true
		}
	}
	return false
}

// GetStats aggregates event data based on grouping and aggregation criteria
// Supports query parameters:
// - groupBy: Field to group results by
// - aggregates: Aggregation operation (count, sum, avg)
// - filters: JSON string for complex MongoDB queries
// - Any other query parameter will be used as a filter (for backward compatibility)
func GetStats(c *fiber.Ctx) error {
	ctx := context.Background()

	// Extract query parameters
	groupBy := c.Query(constants.ParamGroupBy, "")
	aggregates := c.Query(constants.ParamAggregates, constants.DefaultAggregates)

	// Validate parameters
	if groupBy == "" {
		return httpx.SendResponse(c, httpx.BadRequest("groupBy parameter is required", nil))
	}
	if !isValidAggregation(aggregates) {
		return httpx.SendResponse(c, httpx.BadRequest("Invalid aggregation type", nil))
	}

	// Parse JSON filters if provided
	var filters bson.M
	if filterStr := c.Query("filters"); filterStr != "" {
		var err error
		filters, err = parseJSONFilters(filterStr)
		if err != nil {
			return httpx.SendResponse(c, httpx.BadRequest("Invalid filters parameter", err))
		}
	} else {
		// Extract simple filters for backward compatibility
		filters = bson.M{}
		for key, values := range c.Queries() {
			if key != constants.ParamGroupBy && key != constants.ParamAggregates {
				filters[key] = values
			}
		}
	}

	// Perform aggregation query
	stats, err := internalUtils.AggregateStats(ctx, filters, groupBy, aggregates)
	if err != nil {
		return httpx.SendResponse(c, httpx.InternalServerError("Failed to fetch stats", err))
	}

	return httpx.SendResponse(c, httpx.OK("Stats retrieved successfully", fiber.Map{
		constants.ParamGroupBy:    groupBy,
		constants.ParamAggregates: aggregates,
		"stats":                   stats,
	}))
}

// GetTimeSeries generates time-based aggregations of event data
// Supports query parameters:
// - interval: Time grouping interval (hour, day, week, month)
// - aggregates: Aggregation operation (count, sum, avg)
// - filters: JSON string for complex MongoDB queries
// - Any other query parameter will be used as a filter (for backward compatibility)
func GetTimeSeries(c *fiber.Ctx) error {
	ctx := context.Background()

	// Extract query parameters
	aggregates := c.Query(constants.ParamAggregates, constants.DefaultAggregates)
	interval := c.Query(constants.ParamInterval, constants.DefaultInterval)

	// Validate parameters
	if !isValidAggregation(aggregates) {
		return httpx.SendResponse(c, httpx.BadRequest("Invalid aggregation type", nil))
	}
	if !isValidTimeInterval(interval) {
		return httpx.SendResponse(c, httpx.BadRequest("Invalid time interval", nil))
	}

	// Parse JSON filters if provided
	var filters bson.M
	if filterStr := c.Query("filters"); filterStr != "" {
		var err error
		filters, err = parseJSONFilters(filterStr)
		if err != nil {
			return httpx.SendResponse(c, httpx.BadRequest("Invalid filters parameter", err))
		}
	} else {
		// Extract simple filters for backward compatibility
		filters = bson.M{}
		for key, values := range c.Queries() {
			if key != constants.ParamAggregates && key != constants.ParamInterval {
				filters[key] = values
			}
		}
	}

	// Perform time-series query
	timeSeries, err := internalUtils.AggregateTimeSeries(ctx, filters, interval, aggregates)
	if err != nil {
		return httpx.SendResponse(c, httpx.InternalServerError("Failed to fetch time series", err))
	}

	return httpx.SendResponse(c, httpx.OK("Time series retrieved successfully", fiber.Map{
		constants.ParamInterval:   interval,
		constants.ParamAggregates: aggregates,
		"timeSeries":              timeSeries,
	}))
}

// parseJSONFilters parses a JSON string into MongoDB filters
func parseJSONFilters(filterStr string) (bson.M, error) {
	var filters bson.M
	if err := json.Unmarshal([]byte(filterStr), &filters); err != nil {
		return nil, fmt.Errorf("invalid filters JSON: %w", err)
	}
	return filters, nil
}
