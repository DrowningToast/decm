package certificate_share_handler

import (
	"sync"
	"time"
)

type cacheEntry struct {
	data      []byte
	expiresAt time.Time
}

type responseCache struct {
	store sync.Map
}

func (c *responseCache) Get(key string) ([]byte, bool) {
	v, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}
	entry := v.(cacheEntry)
	if time.Now().After(entry.expiresAt) {
		c.store.Delete(key)
		return nil, false
	}
	return entry.data, true
}

func (c *responseCache) Set(key string, data []byte, ttl time.Duration) {
	c.store.Store(key, cacheEntry{data: data, expiresAt: time.Now().Add(ttl)})
}

func (c *responseCache) Delete(key string) {
	c.store.Delete(key)
}
