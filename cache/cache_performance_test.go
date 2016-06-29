package cache

import (
	"testing"
)

func BenchmarkCacheSetStringTest(b *testing.B) {
	cache := NewCache()
	presentTime = &PresentTime{0}
	// run the Set function b.N times
	for n := 0; n < b.N; n++ {
		cache.Set("key", "value", presentTime)
	}
}

func BenchmarkCacheSetListTest(b *testing.B) {
	cache := NewCache()
	presentTime = &PresentTime{0}
	value := [2]string{"Penn", "Teller"}
	// run the Set function b.N times
	for n := 0; n < b.N; n++ {
		cache.Set("key", value, presentTime)
	}
}

func BenchmarkCacheSetMapTest(b *testing.B) {
	cache := NewCache()
	presentTime = &PresentTime{0}
	value := map[string]int{"rsc": 3711, "r": 2138}
	// run the Set function b.N times
	for n := 0; n < b.N; n++ {
		cache.Set("key", value, presentTime)
	}
}

func BenchmarkCacheParallel(b *testing.B) {
	cache := NewCache()
	b.RunParallel(func(pb *testing.PB) {
		presentTime = &PresentTime{0}
		for pb.Next() {
			cache.Set("key", "value", presentTime)
		}
	})
}

func BenchmarkCacheGetValueTest(b *testing.B) {
	cache := NewCache()
	presentTime = &PresentTime{0}
	value := "value"
	cache.Set("key", value, presentTime)
	for n := 0; n < b.N; n++ {
		cache.Get("key")
	}
}

func BenchmarkCacheUpdateValueTest(b *testing.B) {
	cache := NewCache()
	presentTime = &PresentTime{0}
	value := "New value"
	cache.Set("key", value, presentTime)
	for n := 0; n < b.N; n++ {
		cache.Upd("key", value)
	}
}

func BenchmarkCacheDeleteValueTest(b *testing.B) {
	cache := NewCache()
	presentTime = &PresentTime{0}
	for i := 0; i < 10000; i++ {
		cache.Set("key", i, presentTime)
	}
	for n := 0; n < b.N; n++ {
		cache.Del("key")
	}
}

func BenchmarkCacheKeysValueTest(b *testing.B) {
	cache := NewCache()
	presentTime = &PresentTime{0}
	for i := 0; i < 10000; i++ {
		cache.Set("key", i, presentTime)
	}
	for n := 0; n < b.N; n++ {
		cache.Keys(`\d`)
	}
}
