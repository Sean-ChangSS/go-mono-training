package localcache

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSet(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{"string", "key", "value"},
		{"int", "key", 123},
		{"float", "key", 1.23},
		{"bool", "key", true},
		{"slice", "key", []string{"a", "b", "c"}},
		{"map", "key", map[string]int{"a": 1, "b": 2, "c": 3}},
		{"struct", "key", struct {
			Name string
			Age  int
		}{"bob", 18}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache := newCache()
			cache.Set(tc.key, tc.value)

			diff := cmp.Diff(tc.value, cache.data[tc.key].value)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{"string", "key", "value"},
		{"int", "key", 123},
		{"float", "key", 1.23},
		{"bool", "key", true},
		{"slice", "key", []string{"a", "b", "c"}},
		{"map", "key", map[string]int{"a": 1, "b": 2, "c": 3}},
		{"struct", "key", struct {
			Name string
			Age  int
		}{"bob", 18}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache := newCache()
			cache.data[tc.key] = cacheItem{
				value:       tc.value,
				expireation: time.Now().Add(cacheRetainTime),
			}

			got, ok := cache.Get(tc.key)
			if !ok {
				t.Fatalf("key not found")
			}

			diff := cmp.Diff(tc.value, got)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSetValueWithDuplicateKey(t *testing.T) {
	key := "key"
	value1 := "value1"
	value2 := "value2"
	cache := New()
	cache.Set(key, value1)
	cache.Set(key, value2)
	got, _ := cache.Get(key)
	diff := cmp.Diff(value2, got)
	if diff != "" {
		t.Fatalf(diff)
	}
}

func TestCleanUnusedCacheAfter30Seconds(t *testing.T) {
	key := "key"
	value := "value"
	cache := New()
	cache.Set(key, value)

	time.Sleep(30 * time.Second)

	_, ok := cache.Get(key)
	if ok {
		t.Fatalf("cache is not expired")
	}
}
