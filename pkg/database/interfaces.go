package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseService interface {
	Connect() error
	Disconnect(ctx context.Context) error
	GetCollection(name string) *mongo.Collection
	InsertOne(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error)
}
