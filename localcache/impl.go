package localcache

import "time"

const (
	cacheRetainTime = 30 * time.Second
)

type cache struct {
	data map[string]interface{}
}

func (c *cache) Get(key string) (interface{}, bool) {
	value, ok := c.data[key]
	return value, ok
}

func (c *cache) Set(key string, value interface{}) {
	c.data[key] = value
}

func newCache() *cache {
	return &cache{
		data: make(map[string]interface{}),
	}
}
