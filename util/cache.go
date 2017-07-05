package util

import "sync"

type (
	// Cache implements a simple cache
	Cache interface {
		Get(key Any, res CacheResolver) Any
	}

	// CacheResolver is used to resolve the value of a cached item
	CacheResolver func() Any

	cache struct {
		sync.RWMutex
		data valueMap
	}

	valueMap map[Any]Any
)

// NewCache creates a new synchronous Cache instance
func NewCache() Cache {
	return &cache{
		data: valueMap{},
	}
}

func (c *cache) Get(key Any, res CacheResolver) Any {
	c.RLock()
	if r, ok := c.data[key]; ok {
		c.RUnlock()
		return r
	}
	c.RUnlock()

	c.Lock()
	defer c.Unlock()
	r := res()
	c.data[key] = r
	return r
}
