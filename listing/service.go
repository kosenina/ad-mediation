package listing

import (
	"github.com/kosenina/ad-mediation/models"
	"github.com/kosenina/ad-mediation/objectcache"
)

// Service provides ad network listing operations
type Service interface {
	GetAdNetworkList() (models.AdNetworkList, error)
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
func (s *service) GetAdNetworkList() (models.AdNetworkList, error) {
	// Get ad network list from cache
	// If ad network list not cached get from persistant layer and add it to the cache
	return s.repo.Get()
}
