package zigzag

import (
        "zigzag/cache"
        "sync"
        )

var store *cache.Cache

var once sync.Once

// type Clock struct{
//   ex int64
// }

// func (c *Clock) Now() time.Time { return time.Now() }
// func (c *Clock) Duration() time.Duration { return time.Duration(c.ex) * time.Minute }

func Set(key string, value interface {},  m cache.Momenter) {

  once.Do(func() {
      store = &cache.Cache{
        Items: make(map[string]*cache.Item),
      }
  })

  store.Set(key, value, m)
}

func Get(key string) (interface {}, bool) {
  if item, found := store.Get(key); found {
    return item.Object, true
  }
  return nil, false
}

func Del(key string) {
  store.Del(key)
}

func Upd(key string, value interface {}) bool {
  return store.Upd(key, value)
}

func Keys(pattern string) []string {
  return store.Keys(pattern)
}

func GetCache() *cache.Cache{
  return store
}

func DelRandomExpires(num int) int{
  length := len(store.Items)
  expiresRemoved := 0
  i := 0
  // The Go runtime actually randomizes the map iteration order
  for k, v := range store.Items {
    if (i == length) { return expiresRemoved }
    if (i == num) {
      correlation := float64(expiresRemoved) / float64(i)
      if (correlation < 0.25) {
        return expiresRemoved
      } else {
        num *= 2
      }
    }
    if v.Expired() {
      store.Del(k)
      expiresRemoved +=1
    }
    i+=1
  }
  return expiresRemoved
}

