// Package localcache provides a simple in-memory cache.
package localcache

// Cache provides a simple interface for a local cache.
// The guaranteed lifetime of cache entry is 30s.
type Cache interface {

	// Get returns the value for the given key, and a boolean indicating key exists or not.
	Get(key string) (interface{}, bool)

	// Set sets the value for the given key.
	// If the key already exists, the value is overwritten.
	Set(key string, value interface{})
}

// New returns a new Cache instance.
func New() Cache {
	return newCache()
}
