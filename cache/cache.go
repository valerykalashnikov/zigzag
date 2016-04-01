package cache

import (
        "sync"
        "time"
        "regexp"
        )

type Momenter interface {
  Now() time.Time
  Duration() time.Duration
}

type Item struct {
  Object interface {}
  ExpireAt int64
}

func (i *Item) Expired() bool {
  if (i.ExpireAt == 0) {return false}
  return i.ExpireAt <= time.Now().UnixNano()
}

type Cache struct {
  Items map[string]*Item
  mux   sync.RWMutex
}

func (c *Cache) Set(key string, value interface {}, moment Momenter) {
  var expireAt int64
  duration := moment.Duration()
  if duration != 0 {
    expireAt = moment.Now().Add(duration).UnixNano()
  }
  c.mux.Lock()
  c.Items[key] = &Item{
    Object:   value,
    ExpireAt: expireAt,
  }
  c.mux.Unlock()
}

func (c *Cache) Get(key string) (*Item, bool)  {
  if c == nil { return nil, false }
  c.mux.RLock()
  v, ok := c.Items[key]
  c.mux.RUnlock()
  if ok {
    expired := v.Expired()
    if expired {
      return nil, false
    } else {
      return v, true
    }
  }

  return nil, false
}

func (c *Cache) Upd(key string, newValue interface {}) bool {
  // it saves TTL
  item, found := c.Items[key]
  if !found {return false}
  c.mux.Lock()
  item.Object = newValue
  c.mux.Unlock()
  return true
}

func (c *Cache) Del(key string) {
  c.mux.Lock()
  delete(c.Items, key)
  c.mux.Unlock()
}

func (c *Cache) Keys(pattern string) []string {
  keys := make([]string, 0, len(c.Items))
  var validKey = regexp.MustCompile(pattern)
  for k := range c.Items {
      if validKey.MatchString(k) {
        keys = append(keys, k)
      }
  }
  return keys
}
