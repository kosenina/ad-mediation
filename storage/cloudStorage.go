package storage

import (
	"context"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/datastore"
	"github.com/kosenina/ad-mediation/models"
)

const (
	kindName = "adMediation"
)

var storageClient *datastore.Client = nil

// CloudStorage persist data in Google cloud storage
type CloudStorage struct {
	client *datastore.Client
}

// NewGoogleCloudStorage creates CloudStorage instance (is safe to use concurrently)
func NewGoogleCloudStorage(projectID string) *CloudStorage {
	s := new(CloudStorage)

	if storageClient != nil {
		s.client = storageClient
		return s
	}

	var mux sync.Mutex
	mux.Lock()
	defer mux.Unlock()

	if storageClient != nil {
		s.client = storageClient
		return s
	}

	ctx := context.Background()
	var err error
	storageClient, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("ERROR: Could not create datastore client: %v", err)
	}
	s.client = storageClient
	return s
}

// Get returns AdNetworkList from MongoDB
func (s *CloudStorage) Get(documentID string) (models.AdNetworkList, error) {
	var results []*models.AdNetworkList
	var result models.AdNetworkList

	// Create a query to fetch all Task entities, ordered by "created".
	query := datastore.NewQuery(kindName).Filter("id", documentID)
	keys, err := s.client.GetAll(context.Background(), query, &results)
	if err != nil {
		log.Printf("ERROR: Failed to query document by ID %s, error: %v", documentID, err)
		return result, err
	}

	if len(keys) > 0 {
		return *results[0], nil
	}
	return result, fmt.Errorf("Document with ID %s does not exists", documentID)
}

// Upsert inserts or updates AdNetworkList in MongoDB
func (s *CloudStorage) Upsert(data models.AdNetworkList) error {
	var results []*models.AdNetworkList
	ctx := context.Background()

	// Create a query to fetch document by ID.
	query := datastore.NewQuery(kindName).Filter("id", data.ID)
	keys, err := s.client.GetAll(ctx, query, &results)
	if err != nil {
		log.Printf("ERROR: Failed to upsert document by ID %s, error: %v", data.ID, err)
		return err
	}

	if len(keys) > 0 {
		key := *keys[0]
		idKey := datastore.IDKey(kindName, key.ID, nil)

		// In a transaction load each task, set done to true and store.
		_, err := s.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			if _, err := tx.Put(idKey, data); err != nil {
				return err
			}
			return err
		})
		return err
	}

	// Document does not exists, do insert
	key := datastore.IncompleteKey(kindName, nil)
	_, insertionErr := s.client.Put(ctx, key, data)

	return insertionErr
}

// Ping checks if Cloud storage is available
func (s *CloudStorage) Ping() error {
	return nil
}
