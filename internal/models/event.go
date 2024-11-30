package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	Id         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name       string                 `bson:"name" json:"name" validate:"required"`
	Properties map[string]interface{} `bson:"properties" json:"properties"`
	CreatedAt  time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time              `bson:"updated_at" json:"updated_at"`
}
