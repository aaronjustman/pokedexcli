package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	interval time.Duration
	mutex    *sync.Mutex
	ticker   *time.Ticker
}

type cacheEntry struct {
	createdAt time.Time
	value     []byte
}

func NewCache(interval time.Duration) Cache {
	newCache := Cache{
		interval: interval,
		ticker:   time.NewTicker(interval),
	}
	newCache.reapCache()

	return newCache
}

func (c *Cache) Add(key string, value []byte) {
	c.mutex.Lock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		value:     value,
	}
	c.mutex.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	entry, ok := c.cache[key]
	if !ok {
		c.mutex.Unlock()
		return nil, false
	}

	c.mutex.Unlock()
	return entry.value, ok
}

func (c *Cache) reapCache() {
	c.mutex.Lock()
	for {
		select {
		case t := <-c.ticker.C:
			for key := range c.cache {
				value, _ := c.cache[key]
				if t.Sub(value.createdAt) > 0 {
					delete(c.cache, key)
					c.mutex.Unlock()
				}
			}
		}
	}
}
