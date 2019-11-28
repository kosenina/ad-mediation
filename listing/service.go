package listing

import (
	"fmt"
	"log"

	"github.com/kosenina/ad-mediation/models"
	"github.com/kosenina/ad-mediation/objectcache"
)

// Service provides ad network listing operations
type Service interface {
	GetAdNetworkList(string) (models.AdNetworkList, error)
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
func (s *service) GetAdNetworkList(documentID string) (models.AdNetworkList, error) {
	// Get ad network list from cache
	obj := s.cache.Get(fmt.Sprintf("obj:%s", documentID))
	if obj != nil {
		adNetworkList, ok := obj.(models.AdNetworkList)
		if ok {
			log.Println("INFO: got adNetworkList from cache.")
			return adNetworkList, nil
		}
		log.Println("ERROR: Failed to cast cached object to AdNetworkList struct.")
	}

	// If ad network list not cached get from persistant layer and add it to the cache
	return s.repo.Get(documentID)
}
