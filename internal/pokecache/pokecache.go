package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]cacheEntry
	mu       sync.RWMutex
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

func (c *Cache) Get(url string) ([]struct {
	Name string "json:\"name\""
}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.cacheMap[url]

	return val.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)

	for {
		c.mu.Lock()
		for k, v := range c.cacheMap {
			if time.Since(v.createdAt) > interval {
				delete(c.cacheMap, k)
			}
		}
		c.mu.Unlock()

		<-ticker.C
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}

	go cache.reapLoop(interval)

	return cache
}
