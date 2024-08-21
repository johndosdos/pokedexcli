package pokecache

import (
	"sync"
	"time"
)

type cache struct {
	cacheMap map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *cache) Add(url string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheMap[url] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *cache) Get(url string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.cacheMap[url]

	return val.val, ok
}

func NewCache(interval time.Duration) *cache {
	return &cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}
}
