package services

import (
	"context"
	"events-api/internal/models"
	"events-api/pkg/database"
	"events-api/pkg/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateEvent(ctx context.Context, properties map[string]interface{}) (*models.Event, error) {
	event := models.Event{
		Id:         primitive.NewObjectID(),
		Properties: properties,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result, err := database.DBClient.Database(utils.GetEnv("DB_NAME")).Collection("events").InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	event.Id = result.InsertedID.(primitive.ObjectID)
	return &event, nil
}
