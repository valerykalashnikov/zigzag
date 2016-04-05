package cache

import (
        "testing"
        )

func BenchmarkSetStringTest(b *testing.B) {
  setup()
  presentTime = &PresentTime{0}
  // run the Set function b.N times
  for n := 0; n < b.N; n++ {
    cache.Set("key", "value", presentTime)
  }
}

func BenchmarkSetListTest(b *testing.B) {
  setup()
  presentTime = &PresentTime{0}
  value := [2]string{"Penn", "Teller"}
  // run the Set function b.N times
  for n := 0; n < b.N; n++ {
    cache.Set("key", value, presentTime)
  }
}

func BenchmarkSetMapTest(b *testing.B) {
  setup()
  presentTime = &PresentTime{0}
  value := map[string]int{"rsc": 3711, "r": 2138}
  // run the Set function b.N times
  for n := 0; n < b.N; n++ {
    cache.Set("key", value, presentTime)
  }
}


func BenchmarkGetValueTest(b *testing.B) {
  setup()
  presentTime = &PresentTime{0}
  value := "value"
  cache.Set("key", value, presentTime)
  for n := 0; n < b.N; n++ {
    cache.Get("key")
  }
}


func BenchmarkUpdateValueTest(b *testing.B) {
  setup()
  presentTime = &PresentTime{0}
  value := "New value"
  cache.Set("key", value, presentTime)
  for n := 0; n < b.N; n++ {
    cache.Upd("key", value)
  }
}

func BenchmarkDeleteValueTest(b *testing.B) {
  setup()
  presentTime = &PresentTime{0}
  for i := 0; i < 10000; i++ {
    cache.Set("key", i, presentTime)
  }
  for n := 0; n < b.N; n++ {
    cache.Del("key")
  }
}

func BenchmarkKeysValueTest(b *testing.B) {
  setup()
  presentTime = &PresentTime{0}
  for i := 0; i < 10000; i++ {
    cache.Set("key", i, presentTime)
  }
  for n := 0; n < b.N; n++ {
    cache.Keys(`\d`)
  }
}
