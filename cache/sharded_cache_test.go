package cache

import (
          "testing"
          "reflect"
          "time"
        )

func TestShardedCacheSet(t *testing.T) {
  presentTime = &PresentTime{0}
  shardedCache := NewShardedCache()
  // when TTL is not defined
  var testsWithoutTTL = []struct {
    key    string
    value    interface {}
    expected *Item
  }{
    {
      "key",
      "value",
      &Item{"value", 0},
    },
    {
      "key",
      [2]string{"Penn", "Teller"},
      &Item{[2]string{"Penn", "Teller"}, 0},
    },
    {
      "key",
      map[string]int{"rsc": 3711, "r": 2138},
      &Item{map[string]int{"rsc": 3711, "r": 2138}, 0},
    },
  }
  for _, tt := range testsWithoutTTL {
    shardedCache.Set(tt.key, tt.value, presentTime)
    shard := shardedCache.GetShard(tt.key)

    actual := shard.items[tt.key]
    if !reflect.DeepEqual(actual, tt.expected) {
      t.Errorf("Set: expected %v, actual %v", tt.expected, actual)
    }
  }

  ex := int64(30)
  key :=   "key"
  value := "value"
  duration := time.Duration(ex) * time.Minute
  ancientTime = &AncientTime{duration}
  shardedCache.Set(key, value, ancientTime)
  shard := shardedCache.GetShard(key)

  actual := shard.items[key]
  if actual.ExpireAt == 0 {
    t.Errorf("Set: ExpireAt should be greater than nil")
  }
  if actual.ExpireAt != ancientTime.Now().Add(duration).UnixNano() {
    t.Errorf("Set: ExpireAt should be now + duration")
  }
}

func TestShardedCacheGet(t *testing.T) {
  shardedCache := NewShardedCache()
  // when ttl is not defined
  presentTime = &PresentTime{0}
  key := "key"
  expected := &Item{"value", 0}

  shardedCache.Set(key, "value", presentTime)
  actual, _ := shardedCache.Get(key)

  if (*actual != *expected) {
    t.Errorf("Get: expected %v, actual %v", expected, actual)
  }
  //when key is not present in storage
  shard := shardedCache.GetShard(key)

  delete(shard.items, "key")

  actual, found := shardedCache.Get(key)
  if (found != false) {
    t.Error("Get: expected false while returning empty value, found flag is", found)
  }

  //when key is presented in storage but it is outdated
  ancientTime = &AncientTime{10}
  delete(shard.items, "key")
  shardedCache.Set(key, "value", ancientTime)
  actual, found = shardedCache.Get(key)
  if (found == true) {
    t.Error("Get: Expired value shouldn't be returned, actual =", found)
  }
}

func TestShardedCacheUpd(t *testing.T) {
  shardedCache := NewShardedCache()

  //it shouldn't change expiration
  key :=   "key"
  value := "value"
  ex := int64(30)

  duration := time.Duration(ex) * time.Minute
  ancientTime = &AncientTime{duration}

  shardedCache.Set(key, value, ancientTime)
  expectedValue := "newValue"
  shardedCache.Upd(key, expectedValue)
  shard := shardedCache.GetShard(key)

  actualValue := shard.items[key].Object

  if (actualValue != expectedValue) {
    t.Errorf("Upd: Object should be changed, expect %v, actual %v", expectedValue, actualValue)
  }

}

func TestShardedCacheDel(t *testing.T) {
  shardedCache := NewShardedCache()

  key := "key"

  expected := false

  presentTime = &PresentTime{0}
  shardedCache.Set(key, "value", presentTime)

  shardedCache.Del(key)

  _, found := shardedCache.Get(key)

  if (found != expected) {
    t.Errorf("Get: expected %v, actual %v", expected, found)
  }

}

func TestShardedCacheKeys(t *testing.T) {
  shardedCache := NewShardedCache()

  presentTime = &PresentTime{0}
  shardedCache.Set("adam[23]", "value", presentTime)
  shardedCache.Set("eve[7]", "value", presentTime)
  shardedCache.Set("Job[48]", "value", presentTime)
  shardedCache.Set("snakey", "value", presentTime)

  pattern := "^[a-z]+[[0-9]+]$"

  keys := shardedCache.Keys(pattern)

  if (!stringInSlice("adam[23]", keys)) {
    t.Errorf("Keys: expect %v, got %v", "adam[23]", keys[0])
  }
  if (!stringInSlice("eve[7]", keys)) {
    t.Errorf("Keys: expect %v, got %v", "eve[7]", keys[1])
  }

  if (len(keys) != 2) {
    t.Errorf("Keys: length of keys should be 2, got: %v", len(keys))
  }
}
