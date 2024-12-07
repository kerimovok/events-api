package database

import (
	"context"
	"events-api/pkg/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBClient *mongo.Client

const (
	defaultTimeout = 10 * time.Second
)

func ConnectDB() error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	uri := utils.GetEnv("DB_URI")
	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection with timeout
	var result bson.M
	if err := client.Database(utils.GetEnv("DB_NAME")).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	DBClient = client
	utils.LogInfo("Successfully connected to MongoDB!")
	return nil
}

func DisconnectDB(ctx context.Context) error {
	if DBClient == nil {
		return nil
	}

	// Create context with timeout if none provided
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
	}

	if err := DBClient.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	return nil
}
