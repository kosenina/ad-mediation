package adding

import (
	"fmt"
	"log"
	"time"

	"github.com/kosenina/ad-mediation/models"
	"github.com/kosenina/ad-mediation/objectcache"
	"github.com/kosenina/ad-mediation/utils"
)

// Service provides ad network listing operations
type Service interface {
	UpsertAdNetworkList(models.AdNetworkList) error
}

type service struct {
	repo  models.Repository
	cache objectcache.ObjectCache
}

// NewService creates service with the necessary dependencies
func NewService(repo models.Repository, cache objectcache.ObjectCache) Service {
	return &service{repo, cache}
}

// GetAdNetworkList returns AdNetworkList struct
func (s *service) UpsertAdNetworkList(data models.AdNetworkList) error {
	// Populate necessary data
	data.Created = time.Now()
	data.ID = utils.GetAdNetworkListID(data.Created)

	// try to save document
	result := s.repo.Upsert(data)
	if result == nil {
		s.cache.Remove("obj:adNetworkList")
	} else {
		log.Printf(fmt.Sprintf("ERROR: Failed to upsert ad network list with ID: %s, internal error: %v", data.ID, result))
	}
	return result
}
