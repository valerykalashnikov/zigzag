package cache

import (
          "fmt"
          "sync"
          "hash/fnv"
          "regexp"
        )

type ShardedCache map[string]*CacheShard

type CacheShard struct {
  items map[string]*Item
  mux *sync.RWMutex
}

func (c ShardedCache) Set(key string, value interface {}, moment Momenter) {
  var expireAt int64
  duration := moment.Duration()
  if duration != 0 {
    expireAt = moment.Now().Add(duration).UnixNano()
  }
  shard := c.getShard(key)
  shard.mux.Lock()
  defer shard.mux.Unlock()
  shard.items[key] = &Item{
    Object:   value,
    ExpireAt: expireAt,
  }
}

func (c ShardedCache) Get(key string) (*Item, bool) {
  shard := c.getShard(key)
  shard.mux.RLock()
  defer shard.mux.RUnlock()

  if item, ok := shard.items[key]; ok { return item, true }
  return nil, false
}

func (c ShardedCache) Upd(key string, newValue interface {}) bool {
  // it saves TTL
  shard := c.getShard(key)
  shard.mux.Lock()
  defer shard.mux.Unlock()
  item, found := shard.items[key]
  if !found {return false}
  item.Object = newValue
  return true
}

func (c ShardedCache) Del(key string) {
  shard := c.getShard(key)
  shard.mux.Lock()
  defer shard.mux.Unlock()
  delete(shard.items, key)
}

func (c ShardedCache) Keys(pattern string) []string{
  var keys []string
  var validKey = regexp.MustCompile(pattern)
  for _, shard := range c {

    for k := range shard.items {
      if validKey.MatchString(k) {
        keys = append(keys, k)
      }
    }
  }

  return keys
}

func (c ShardedCache) Items() map[string]*Item {
  results:=  make(map[string]*Item)
  for _, shard := range c {
    for key, item := range shard.items {
      results[key] = item
    }
  }
  return results
}

func (c ShardedCache) getShard(key string) *CacheShard {
  hasher := fnv.New64()
  hasher.Write([]byte(key))
  shardKey :=  fmt.Sprintf("%x", hasher.Sum(nil))[0:2]
  return c[shardKey]
}


func NewShardedCache() DataStore {
  c := make(ShardedCache, 256)
  for i := 0; i < 256; i++ {
    c[fmt.Sprintf("%02x", i)] = &CacheShard{
      items: make(map[string]*Item, 2048),
      mux: new(sync.RWMutex),
    }
  }
  return c
}
