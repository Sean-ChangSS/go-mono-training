package localcache

import (
	"runtime"
	"sync"
	"time"
)

const (
	cacheRetainTime = 30 * time.Second
)

type cachePointer struct {
	*cache
}

type cacheItem struct {
	value       interface{}
	expireation time.Time
}

type janitor struct {
	stop chan bool
}

type cache struct {
	data    map[string]cacheItem
	mu      sync.RWMutex
	janitor janitor
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

func stopClean(c *cachePointer) {
	c.janitor.stop <- true
}

func newCache() *cache {
	c := &cache{
		data: make(map[string]cacheItem),
	}

	// Start a background goroutine to periodically clean the cache
	C := &cachePointer{c}
	go c.clean()
	runtime.SetFinalizer(C, stopClean)

	return c
}
