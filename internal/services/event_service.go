package services

import (
	"context"
	"events-api/internal/models"
	"events-api/pkg/database"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type eventService struct {
	db database.DatabaseService
}

func NewEventService(db database.DatabaseService) EventService {
	return &eventService{
		db: db,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, name string, properties map[string]interface{}) (*models.Event, error) {
	event := models.Event{
		Id:         primitive.NewObjectID(),
		Name:       name,
		Properties: properties,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result, err := s.db.InsertOne(ctx, "events", event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	event.Id = result.InsertedID.(primitive.ObjectID)
	return &event, nil
}
