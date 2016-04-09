package cache

import (
        "testing"
        )

func TestCreateCache(t *testing.T) {
  cache, _ := CreateCache()
  _, ok := cache.(*Cache)
  if !ok {
    t.Errorf("Cache test: returned value shoulbe be of 'Cache' type")
  }
}


func TestCreateShardedCache(t *testing.T) {
  shardedCache, _ := CreateCache("sharded")
  _, ok := shardedCache.(ShardedCache)
  if !ok {
    t.Errorf("Cache test: returned value shoulbe be of 'ShardedCache' type")
  }
}
