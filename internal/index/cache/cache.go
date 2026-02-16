package cache

import (
	"sync"
	"time"
)

// Entry represents a cached item with expiration
type entry struct {
	value      interface{}
	expiration time.Time
}

// Cache provides a simple in-memory cache with TTL
type Cache struct {
	mu      sync.RWMutex
	items   map[string]*entry
	ttl     time.Duration
	janitor *time.Ticker
	stopCh  chan bool
}

// New creates a new cache with the specified TTL
func New(ttl time.Duration) *Cache {
	c := &Cache{
		items:  make(map[string]*entry),
		ttl:    ttl,
		stopCh: make(chan bool),
	}

	// Start janitor to clean up expired items
	if ttl > 0 {
		c.janitor = time.NewTicker(ttl)
		go c.cleanup()
	}

	return c
}

// Set stores a value in the cache
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &entry{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*entry)
}

// cleanup removes expired items periodically
func (c *Cache) cleanup() {
	for {
		select {
		case <-c.janitor.C:
			c.deleteExpired()
		case <-c.stopCh:
			return
		}
	}
}

// deleteExpired removes all expired items
func (c *Cache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.expiration) {
			delete(c.items, key)
		}
	}
}

// Stop stops the janitor and prevents further cleanup
func (c *Cache) Stop() {
	if c.janitor != nil {
		c.janitor.Stop()
		c.stopCh <- true
	}
}

// Size returns the number of items in the cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
