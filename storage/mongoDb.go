package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

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

var mongoDB *mongo.Database = nil

// MongoDBStorage stores data in Mongo DB
type MongoDBStorage struct {
	db *mongo.Database
}

// NewMongoDBStorage returns a MongoDB storage implementation
func NewMongoDBStorage() (*MongoDBStorage, error) {
	var err error

	s := new(MongoDBStorage)

	if mongoDB != nil {
		s.db = mongoDB
	} else {
		s.db, err = getMongoDatabase()
		if err != nil {
			log.Fatal("ERROR: Failed to connect to MongoDB", err)
		}
	}
	return s, nil
}

// Creates MongoDB client (is safe to use concurrently)
func getMongoDatabase() (*mongo.Database, error) {
	var mux sync.Mutex
	mux.Lock()
	if mongoDB != nil {
		return mongoDB, nil
	}

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

	// Create index to enable query
	mod := mongo.IndexModel{
		Keys: bson.M{
			"created": 1, // index in ascending order
		}, Options: nil,
	}
	col := client.Database(databaseName).Collection(collectionName)
	ind, err := col.Indexes().CreateOne(context.Background(), mod)
	if err != nil {
		log.Printf("ERROR: Failed to create MongoDB index, mongo error: %s", err.Error())
		os.Exit(1) // exit in case of error
	} else {
		log.Printf("INFO: Created MongoDB index: %s", ind)
	}
	mongoDB = client.Database(databaseName)
	mux.Unlock()
	return mongoDB, nil
}

// Get returns AdNetworkList from MongoDB
func (repo *MongoDBStorage) Get(documentID string) (models.AdNetworkList, error) {
	var result models.AdNetworkList
	filter := bson.D{primitive.E{Key: "_id", Value: documentID}}

	collection := repo.db.Collection(collectionName)

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Printf("ERROR: Failed to get document by ID %s, Mongo error: %s.", documentID, err.Error())
		return result, fmt.Errorf("Document with ID %s does not exists", documentID)
	}

	return result, nil
}

// Upsert inserts or updates AdNetworkList in MongoDB
func (repo *MongoDBStorage) Upsert(data models.AdNetworkList) error {
	filter := bson.D{primitive.E{Key: "_id", Value: data.ID}}
	opts := options.Replace().SetUpsert(true)
	collection := repo.db.Collection(collectionName)
	result, err := collection.ReplaceOne(context.TODO(), filter, data, opts)
	if err != nil {
		log.Printf("ERROR: Failed to update document by ID %s, Mongo error: %s", data.ID, err.Error())
		return fmt.Errorf("AdNetworkList with ID %s was not upserted", data.ID)
	}
	log.Printf(
		"INFO: Matched %v documents with ID %s and upserted %v documents, modified %v documents.\n",
		result.MatchedCount,
		data.ID,
		result.UpsertedCount,
		result.ModifiedCount)
	if result.UpsertedCount > 0 || result.ModifiedCount > 0 {
		return nil
	}
	return fmt.Errorf("AdNetworkList with ID %s was not upserted", data.ID)
}

// Ping checks if MongoDB is up and runing
func (repo *MongoDBStorage) Ping() error {
	err := repo.db.Client().Ping(context.TODO(), nil)
	return err
}
