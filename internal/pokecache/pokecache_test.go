package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestCache_AddGet(t *testing.T) {
	cache := NewCache(5 * time.Second)
	defer cache.Stop()

	// Test Add and Get
	url := "https://example.com"
	val := []byte("test data")

	cache.Add(url, val)
	retrievedVal, ok := cache.Get(url)

	if !ok {
		t.Errorf("Value not found in cache for URL: %s", url)
	}
	if string(retrievedVal) != string(val) {
		t.Errorf("Retrieved value does not match stored value.\nExpected: %s\nGot: %s", string(val), string(retrievedVal))
	}
}

func TestCache_ReapLoop(t *testing.T) {
	cache := NewCache(5 * time.Second)
	defer cache.Stop()

	// Test automatic reaping
	url1 := "https://example1.com"
	url2 := "https://example2.com"
	cache.Add(url1, []byte("test data 1"))
	cache.Add(url2, []byte("test data 2"))

	time.Sleep(6 * time.Second) // Wait longer than interval

	_, ok1 := cache.Get(url1)
	_, ok2 := cache.Get(url2)

	if ok1 || ok2 {
		t.Error("Cache entries not reaped as expected.")
	}
}

func TestCache_Stop(t *testing.T) {
	cache := NewCache(5 * time.Second)

	// Attempt to Stop a stopped cache (expect no issues)
	cache.Stop()

	// Basic Add/Get should still function correctly
	url := "testurl"
	data := []byte("testdata")
	cache.Add(url, data)
	if _, ok := cache.Get(url); !ok {
		t.Error("Get failed after Stop - cache should still function for existing entries.")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(5 * time.Second)
	defer cache.Stop()

	numGoroutines := 10
	iterations := 100

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < iterations; j++ {
				url := fmt.Sprintf("https://example%d-%d.com", id, j)
				data := []byte(fmt.Sprintf("test data %d-%d", id, j))

				cache.Add(url, data)
				if _, ok := cache.Get(url); !ok {
					t.Errorf("Concurrent access issue: Value not found for %s", url)
				}
			}
		}(i)
	}

	// Wait to allow goroutines to complete
	time.Sleep(time.Second)
}
