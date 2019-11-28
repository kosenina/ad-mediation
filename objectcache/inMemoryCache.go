package objectcache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// ObjectCache provides in memory cache.
type ObjectCache interface {
	// Get returns cached object.
	Get(string) interface{}
	// Add cache the provided object.
	Add(string, interface{})
	// Remove removes the cached object by key
	Remove(string)
}

// InMemoryCache type implements ObjectCache
type InMemoryCache struct {
	c *cache.Cache
}

// NewInMemoryCache creates in memory cache instance
func NewInMemoryCache() *InMemoryCache {
	s := new(InMemoryCache)
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	s.c = cache.New(5*time.Minute, 10*time.Minute)
	return s
}

// Get provides fetching cached object by key
func (memoryCache *InMemoryCache) Get(key string) interface{} {
	data, found := memoryCache.c.Get(key)
	if found {
		return data
	}
	return nil
}

// Add provides caching object by key
func (memoryCache *InMemoryCache) Add(key string, data interface{}) {
	memoryCache.c.Set(key, data, cache.DefaultExpiration)
}

// Remove implements removing cached object by key
func (memoryCache *InMemoryCache) Remove(key string) {
	memoryCache.c.Delete(key)
}
