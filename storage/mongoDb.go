package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/kosenina/ad-mediation/models"
	"github.com/kosenina/ad-mediation/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURL       = "localhost:27017"
	databaseName   = "adMediation"
	collectionName = "adNetworks"

	adNetworksListID = "outfit7AdNetworks"
)

// MongoDBStorage stores data in Mongo DB
type MongoDBStorage struct {
	db *mongo.Database
}

// NewMongoDBStorage returns a MongoDB storage implementation
func NewMongoDBStorage() (*MongoDBStorage, error) {
	var err error

	s := new(MongoDBStorage)

	s.db, err = getMongoDatabase()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}

	return s, nil
}

// Creates MongoDB client
func getMongoDatabase() (*mongo.Database, error) {
	uri := fmt.Sprintf("mongodb://%s", utils.GetEnv("MONGO_URL", mongoURL))
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}
	return client.Database(databaseName), nil
}

// Get returns AdNetworkList from MongoDB
func (repo *MongoDBStorage) Get() (models.AdNetworkList, error) {
	var result models.AdNetworkList
	filter := bson.D{primitive.E{Key: "_id", Value: adNetworksListID}}

	collection := repo.db.Collection(collectionName)

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println("Failed to get document by ID.", err)
		return result, models.ErrNotFound
	}

	return result, nil
}

// Upsert inserts or updates AdNetworkList in MongoDB
func (repo *MongoDBStorage) Upsert(data models.AdNetworkList) error {
	filter := bson.D{primitive.E{Key: "_id", Value: adNetworksListID}}
	opts := options.Replace().SetUpsert(true)
	collection := repo.db.Collection(collectionName)
	result, err := collection.ReplaceOne(context.TODO(), filter, data, opts)
	if err != nil {
		log.Println("Failed to update document by ID", err)
		return models.ErrUpsertFailed
	}
	fmt.Printf("Matched %v documents and upserted %v documents, modified %v documents.\n", result.MatchedCount, result.UpsertedCount, result.ModifiedCount)
	if result.UpsertedCount > 0 || result.ModifiedCount > 0 {
		return nil
	}
	return models.ErrUpsertFailed
}
