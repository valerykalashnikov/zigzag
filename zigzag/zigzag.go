package zigzag

import (
        "zigzag/cache"
        "sync"
        "time"
        )

var store *cache.Cache

var once sync.Once

type Clock struct{
  duration time.Duration
}

func (c *Clock) Now() time.Time { return time.Now() }
func (c *Clock) Duration() time.Duration { return c.duration }

func Set(key string, value interface {},  optional ...int64 ) {
  var ex int64 = 0

  if len(optional) > 0 {
    ex = optional[0]
  }

  once.Do(func() {
      store = &cache.Cache{
        Items: make(map[string]*cache.Item),
      }
  })

  duration := time.Duration(ex) * time.Minute
  clock := &Clock{duration}
  store.Set(key, value, clock)
}

func Get(key string) (interface {}, bool) {
  item, found := store.Get(key)
  return item, found
}

func Del(key string) {
  store.Del(key)
}

func Upd(key string, value interface {}) bool {
  return store.Upd(key, value)
}


func GetCache() *cache.Cache{
  return store
}

