package database

import (
	"context"
	"events-api/pkg/config"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoService struct {
	client   *mongo.Client
	database string
}

func NewMongoService() DatabaseService {
	return &mongoService{
		database: config.AppConfig.MongoDB.Database,
	}
}

func (m *mongoService) Connect() error {
	uri := config.AppConfig.MongoDB.URI
	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	// Verify connection
	var result bson.M
	if err := client.Database(m.database).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return err
	}

	m.client = client
	fmt.Println("Successfully connected to MongoDB!")
	return nil
}

func (m *mongoService) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *mongoService) GetCollection(name string) *mongo.Collection {
	return m.client.Database(m.database).Collection(name)
}

func (m *mongoService) InsertOne(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	return m.GetCollection(collection).InsertOne(ctx, document)
}
