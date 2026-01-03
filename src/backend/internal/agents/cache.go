package agents

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

// DefaultCacheTTL is the default time-to-live for cached responses
const DefaultCacheTTL = 1 * time.Hour

// DefaultCacheSize is the default maximum number of cached entries
const DefaultCacheSize = 1000

// Cache defines the interface for caching agent responses
type Cache interface {
	// Get retrieves a cached value by key
	Get(key string) ([]byte, bool)

	// Set stores a value with the given TTL
	Set(key string, value []byte, ttl time.Duration)

	// Delete removes a cached value
	Delete(key string)

	// Clear removes all cached values
	Clear()
}

// cacheEntry holds a cached value with its expiration time
type cacheEntry struct {
	value     []byte
	expiresAt time.Time
}

// LRUCache implements Cache using an LRU eviction policy with TTL
type LRUCache struct {
	cache *lru.Cache[string, cacheEntry]
	mu    sync.RWMutex
}

// NewLRUCache creates a new LRU cache with the given size
func NewLRUCache(size int) (*LRUCache, error) {
	if size <= 0 {
		size = DefaultCacheSize
	}

	cache, err := lru.New[string, cacheEntry](size)
	if err != nil {
		return nil, err
	}

	c := &LRUCache{
		cache: cache,
	}

	// Start background cleanup goroutine
	go c.cleanupLoop()

	return c, nil
}

// Get retrieves a cached value, returning false if not found or expired
func (c *LRUCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cache.Get(key)
	if !ok {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.expiresAt) {
		// Don't remove here to avoid write lock, cleanup will handle it
		return nil, false
	}

	return entry.value, true
}

// Set stores a value with the given TTL
func (c *LRUCache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Add(key, cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	})
}

// Delete removes a cached value
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Remove(key)
}

// Clear removes all cached values
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Purge()
}

// cleanupLoop periodically removes expired entries
func (c *LRUCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

// cleanup removes expired entries
func (c *LRUCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	keys := c.cache.Keys()

	for _, key := range keys {
		if entry, ok := c.cache.Peek(key); ok {
			if now.After(entry.expiresAt) {
				c.cache.Remove(key)
			}
		}
	}
}

// NoOpCache is a cache that does nothing (for testing or when caching is disabled)
type NoOpCache struct{}

// NewNoOpCache creates a new no-op cache
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// Get always returns false
func (c *NoOpCache) Get(_ string) ([]byte, bool) {
	return nil, false
}

// Set does nothing
func (c *NoOpCache) Set(_ string, _ []byte, _ time.Duration) {}

// Delete does nothing
func (c *NoOpCache) Delete(_ string) {}

// Clear does nothing
func (c *NoOpCache) Clear() {}
