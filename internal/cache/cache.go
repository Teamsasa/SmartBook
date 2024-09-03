package cache

import (
	"sync"
	"time"
)

// Cache インターフェースは、キャッシュの基本操作を定義します
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiration time.Duration)
	Delete(key string)
}

// item はキャッシュ内の各アイテムを表します
type item struct {
	value      interface{}
	expiration int64
}

// InMemoryCache は、メモリ内キャッシュの実装です
type InMemoryCache struct {
	items map[string]item
	mu    sync.RWMutex
}

// NewInMemoryCache は新しい InMemoryCache インスタンスを作成します
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		items: make(map[string]item),
	}
}

// Get はキャッシュからアイテムを取得します
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if item.expiration > 0 && item.expiration < time.Now().UnixNano() {
		return nil, false
	}

	return item.value, true
}

// Set はキャッシュにアイテムを設定します
func (c *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp int64
	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	}

	c.items[key] = item{
		value:      value,
		expiration: exp,
	}
}

// Delete はキャッシュからアイテムを削除します
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// StartCleanup は定期的に期限切れのアイテムを削除します
func (c *InMemoryCache) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.deleteExpired()
		}
	}()
}

// deleteExpired は期限切れのアイテムを削除します
func (c *InMemoryCache) deleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.expiration > 0 && v.expiration < now {
			delete(c.items, k)
		}
	}
}
