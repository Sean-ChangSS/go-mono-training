package localcache

import (
	"runtime"
	"sync"
	"time"
)

const (
	cacheRetainTime = 30 * time.Second
)

type cacheItem struct {
	value       interface{}
	expireation time.Time
}

type janitor struct {
	stop chan bool
}

type sharedCache struct {
	data    map[string]cacheItem
	mu      sync.RWMutex
	janitor janitor
}

type cache struct {
	*sharedCache
}

func (c *sharedCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheItem{
		value:       value,
		expireation: time.Now().Add(cacheRetainTime),
	}
}

func (c *sharedCache) Get(key string) (interface{}, bool) {
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

func (c *sharedCache) clean() {
	j := janitor{
		stop: make(chan bool),
	}
	c.janitor = j

	for {
		if <-j.stop {
			break
		}

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

func stopClean(c *cache) {
	c.janitor.stop <- true
}

func newCache() *cache {
	c := &sharedCache{
		data: make(map[string]cacheItem),
	}

	// Start a background goroutine to periodically clean the cache
	C := &cache{c}
	go c.clean()
	runtime.SetFinalizer(C, stopClean)

	return C
}
