package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []struct {
		Name string `json:"name"`
	}
}

func (c *Cache) Add(url string, val []struct {
	Name string `json:"name"`
}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheMap[url] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

	c.mu.Lock()
	defer c.mu.Unlock()
func (c *Cache) Get(url string) ([]struct {
	Name string "json:\"name\""
}, bool) {

	val, ok := c.cacheMap[url]

	return val.val, ok
}

func NewCache(interval time.Duration) *cache {
	return &cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}
}
