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

func TestMultipleSetGetOperation(t *testing.T) {
	key1 := "key1"
	value1 := "value1"
	key2 := "key2"
	value2 := "value2"
	cache := New()
	cache.Set(key1, value1)
	cache.Set(key2, value2)
	got1, _ := cache.Get(key1)
	diff1 := cmp.Diff(value1, got1)
	if diff1 != "" {
		t.Fatalf(diff1)
	}
	got2, _ := cache.Get(key2)
	diff2 := cmp.Diff(value2, got2)
	if diff2 != "" {
		t.Fatalf(diff2)
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
