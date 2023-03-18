package localcache

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSetArbitraryTypeAndGet(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		value  interface{}
		expect interface{}
	}{
		{"string", "key", "value", "value"},
		{"int", "key", 123, 123},
		{"float", "key", 1.23, 1.23},
		{"bool", "key", true, true},
		{"slice", "key", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"map", "key", map[string]int{"a": 1, "b": 2, "c": 3}, map[string]int{"a": 1, "b": 2, "c": 3}},
		{"struct", "key", struct {
			Name string
			Age  int
		}{"bob", 18}, struct {
			Name string
			Age  int
		}{"bob", 18}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache := New()
			cache.Set(tc.key, tc.value)
			got, _ := cache.Get(tc.key)
			diff := cmp.Diff(tc.expect, got)
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
