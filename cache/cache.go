package cache

import (
	"regexp"
	"sync"
	"time"
)

type Momenter interface {
	Now() time.Time
	Duration() time.Duration
}

type Item struct {
	Object   interface{}
	ExpireAt int64
}

func (i *Item) Expired() bool {
	if i.ExpireAt == 0 {
		return false
	}
	return i.ExpireAt <= time.Now().UnixNano()
}

type Cache struct {
	items map[string]*Item
	mux   sync.RWMutex
}

func (c *Cache) Set(key string, value interface{}, moment Momenter) {
	var expireAt int64
	duration := moment.Duration()
	if duration != 0 {
		expireAt = moment.Now().Add(duration).UnixNano()
	}
	c.mux.Lock()
	c.items[key] = &Item{
		Object:   value,
		ExpireAt: expireAt,
	}
	c.mux.Unlock()
}

func (c *Cache) Get(key string) (*Item, bool) {
	if c == nil {
		return nil, false
	}
	c.mux.RLock()
	item, ok := c.items[key]
	c.mux.RUnlock()
	if ok {
		return item, true
	}

	return nil, false
}

func (c *Cache) Upd(key string, newValue interface{}) bool {
	// it saves TTL
	item, found := c.items[key]
	if !found {
		return false
	}
	c.mux.Lock()
	item.Object = newValue
	c.mux.Unlock()
	return true
}

func (c *Cache) Del(key string) {
	c.mux.Lock()
	delete(c.items, key)
	c.mux.Unlock()
}

func (c *Cache) Keys(pattern string) []string {
	keys := make([]string, 0, len(c.items))
	var validKey = regexp.MustCompile(pattern)
	for k := range c.items {
		if validKey.MatchString(k) {
			keys = append(keys, k)
		}
	}
	return keys
}

func (c *Cache) Items() map[string]*Item {
	return c.items
}

func NewCache() DataStore {
	return &Cache{
		items: make(map[string]*Item),
	}
}
