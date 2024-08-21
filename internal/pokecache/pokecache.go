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

func NewCache(interval time.Duration) *cache {
	return &cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}
}
