package zigzag

import (
        "zigzag/cache"
        "testing"
        )

func TestSet(t *testing.T) {
  // when TTL is not defined
  key := "key"
  value := "value"

  Set(key, value)

  expected := &cache.Item{"value", 0}

  cache := GetCache()
  actual := cache.Items[key]
  if *actual != *expected {
    t.Errorf("Set: expected %v, actual %v", expected, actual)
  }

  // when TTL is defined
  ex := int64(30)

  Set(key, value, ex)

  actual = cache.Items[key]
  if actual.ExpireAt == 0 {
    t.Errorf("Set: ExpireAt should be greater than nil")
  }

}

func TestGet(t *testing.T) {
  key := "key"

  expected := cache.Item{"value", 0}
  //setup
  Set(key, "value")

  // run
  actual, _ := Get(key)

  if (actual != expected) {
    t.Errorf("Get: expected %v, actual %v", expected, actual)
  }
}



func TestUpd(t *testing.T) {
  key := "key"

  expected := cache.Item{"newValue", 0}
  //setup
  Set(key, "value")

  Upd(key, "newValue")
  // run
  actual, _ := Get(key)

  if (actual != expected) {
    t.Errorf("Get: expected %v, actual %v", expected, actual)
  }
}

func TestDel(t *testing.T) {
  key := "key"

  expected := false
  //setup
  Set(key, "value")

  //run
  Del(key)

  _, found := Get(key)
  if (found != expected) {
    t.Errorf("Get: expected %v, actual %v", expected, found)
  }

}

func TestKeys(t *testing.T) {
  // cache := GetCache()

  Set("adam[23]", "value")

  pattern := "^[a-z]+[[0-9]+]$"

  keys := Keys(pattern)

  if (keys[0] != "adam[23]") {
    t.Errorf("Keys: expect %v, got %v", "adam[23]", keys[0])
  }

}

