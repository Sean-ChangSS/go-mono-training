package localcache

import (
	"sync"
	"time"
)

const (
	cacheRetainTime = 30 * time.Second
)

type cache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value       interface{}
	expireation time.Time
}

func (c *cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheItem{
		value:       value,
		expireation: time.Now().Add(cacheRetainTime),
	}
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(item.expireation) {
		return nil, false
	}

	return item.value, true
}

func (c *cache) clean() {
	for {
		time.Sleep(cacheRetainTime)

		c.mu.Lock()
		for key, item := range c.data {
			if time.Now().After(item.expireation) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

func newCache() *cache {
	c := &cache{
		data: make(map[string]cacheItem),
	}

	// Start a background goroutine to periodically clean the cache
	go c.clean()

	return c
}
