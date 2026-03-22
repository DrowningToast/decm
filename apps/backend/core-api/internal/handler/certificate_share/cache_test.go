package certificate_share_handler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseCache_GetMiss(t *testing.T) {
	var c responseCache
	data, ok := c.Get("missing")
	assert.False(t, ok)
	assert.Nil(t, data)
}

func TestResponseCache_SetAndGet(t *testing.T) {
	var c responseCache
	payload := []byte(`{"hello":"world"}`)
	c.Set("key1", payload, 5*time.Minute)

	data, ok := c.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, payload, data)
}

func TestResponseCache_GetExpiredEntry(t *testing.T) {
	var c responseCache
	c.Set("key1", []byte("data"), 1*time.Millisecond)

	time.Sleep(5 * time.Millisecond)

	data, ok := c.Get("key1")
	assert.False(t, ok)
	assert.Nil(t, data)
}

func TestResponseCache_Delete(t *testing.T) {
	var c responseCache
	c.Set("key1", []byte("data"), 5*time.Minute)
	c.Delete("key1")

	data, ok := c.Get("key1")
	assert.False(t, ok)
	assert.Nil(t, data)
}

func TestResponseCache_StartCleanup_EvictsExpiredEntries(t *testing.T) {
	var c responseCache
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.Set("expired", []byte("old"), 5*time.Millisecond)
	c.Set("fresh", []byte("new"), 5*time.Minute)

	// Verify both entries are present before cleanup runs
	_, expiredPresent := c.store.Load("expired")
	assert.True(t, expiredPresent, "expired entry should be in the map before cleanup")

	c.StartCleanup(ctx, 10*time.Millisecond)

	// Wait for the entry to expire and for at least one cleanup tick
	time.Sleep(50 * time.Millisecond)

	_, expiredAfter := c.store.Load("expired")
	assert.False(t, expiredAfter, "cleanup goroutine should have evicted the expired entry")

	_, freshPresent := c.store.Load("fresh")
	assert.True(t, freshPresent, "non-expired entry should still be present after cleanup")
}

func TestResponseCache_StartCleanup_StopsOnContextCancel(t *testing.T) {
	var c responseCache
	ctx, cancel := context.WithCancel(context.Background())

	c.StartCleanup(ctx, 10*time.Millisecond)

	// Cancel the context — goroutine should stop without panicking
	cancel()
	time.Sleep(20 * time.Millisecond) // give the goroutine time to exit

	// Nothing to assert — just verifying no panic/deadlock
}
