package cache

import (
        "testing"
        )

func BenchmarkShardedCacheSetStringTest(b *testing.B) {
  cache := NewShardedCache()
  presentTime = &PresentTime{0}
  // run the Set function b.N times
  for n := 0; n < b.N; n++ {
    cache.Set("key", "value", presentTime)
  }
}


func BenchmarkShardedCacheSetListTest(b *testing.B) {
  cache := NewShardedCache()
  presentTime = &PresentTime{0}
  value := [2]string{"Penn", "Teller"}
  // run the Set function b.N times
  for n := 0; n < b.N; n++ {
    cache.Set("key", value, presentTime)
  }
}

func BenchmarkShardedCacheSetMapTest(b *testing.B) {
  cache := NewShardedCache()
  presentTime = &PresentTime{0}
  value := map[string]int{"rsc": 3711, "r": 2138}
  // run the Set function b.N times
  for n := 0; n < b.N; n++ {
    cache.Set("key", value, presentTime)
  }
}

func BenchmarkShardedCacheParallel(b *testing.B) {
  cache := NewShardedCache()
  b.RunParallel(func(pb *testing.PB) {
    presentTime = &PresentTime{0}
    for pb.Next() {
      cache.Set("key", "value", presentTime)
    }
  })
}


func BenchmarkShardedCacheGetValueTest(b *testing.B) {
  cache := NewShardedCache()
  presentTime = &PresentTime{0}
  value := "value"
  cache.Set("key", value, presentTime)
  for n := 0; n < b.N; n++ {
    cache.Get("key")
  }
}


func BenchmarkShardedCacheUpdateValueTest(b *testing.B) {
  cache := NewShardedCache()
  presentTime = &PresentTime{0}
  value := "New value"
  cache.Set("key", value, presentTime)
  for n := 0; n < b.N; n++ {
    cache.Upd("key", value)
  }
}

func BenchmarkShardedCacheDeleteValueTest(b *testing.B) {
  cache := NewShardedCache()
  presentTime = &PresentTime{0}
  for i := 0; i < 10000; i++ {
    cache.Set("key", i, presentTime)
  }
  for n := 0; n < b.N; n++ {
    cache.Del("key")
  }
}

func BenchmarkShardedCacheKeysValueTest(b *testing.B) {
  cache := NewShardedCache()
  presentTime = &PresentTime{0}
  for i := 0; i < 10000; i++ {
    cache.Set("key", i, presentTime)
  }
  for n := 0; n < b.N; n++ {
    cache.Keys(`\d`)
  }
}
