package geecache

import (
	"github.com/AYO-Al/myPackage/geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) Len() int {
	return c.lru.Len()
}

func (c *cache) add(key string, value lru.Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return ByteView{}, false
	}

	if value, ok := c.lru.Get(key); ok {
		return value.(ByteView), true
	}
	return ByteView{}, false
}
