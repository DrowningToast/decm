package certificate_share_handler

import (
	"context"
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

// StartCleanup starts a background goroutine that evicts expired entries every interval.
// It stops when ctx is cancelled.
func (c *responseCache) StartCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				c.store.Range(func(key, value any) bool {
					if value.(cacheEntry).expiresAt.Before(now) {
						c.store.Delete(key)
					}
					return true
				})
			case <-ctx.Done():
				return
			}
		}
	}()
}
