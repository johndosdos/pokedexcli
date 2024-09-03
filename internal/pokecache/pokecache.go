package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	CacheMap map[string]cacheEntry
	mu       sync.RWMutex
	interval time.Duration
	done     chan struct{}
}

type cacheEntry struct {
	CreatedAt time.Time
	val       []byte
}

func (c *Cache) Add(url string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CacheMap[url] = cacheEntry{
		CreatedAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(url string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.CacheMap[url]
	if !ok || time.Since(entry.CreatedAt) > c.interval {
		return nil, false
	}
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for k, v := range c.CacheMap {
				if time.Since(v.CreatedAt) > c.interval {
					delete(c.CacheMap, k)
				}
			}
			c.mu.Unlock()

		case <-c.done:
			return
		}

	}
}

func (c *Cache) Stop() {
	close(c.done)
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		CacheMap: make(map[string]cacheEntry),
		interval: interval,
		done:     make(chan struct{}),
	}
	go cache.reapLoop()
	return cache
}
