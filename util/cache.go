package util

import "sync"

// Cache implements a simple cache
type Cache interface {
	Get(key Any, res CacheResolver) Any
}

type cache struct {
	sync.RWMutex
	data valueMap
}

// CacheResolver is used to resolve the value of a cached item
type CacheResolver func() Any
type valueMap map[Any]Any

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
