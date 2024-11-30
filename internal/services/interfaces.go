package services

import (
	"context"
	"events-api/internal/models"
)

type EventService interface {
	CreateEvent(ctx context.Context, name string, properties map[string]interface{}) (*models.Event, error)
}
