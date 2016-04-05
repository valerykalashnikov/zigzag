package zigzag

import (
        "github.com/valerykalashnikov/zigzag/cache"
        "testing"
        "time"
        )

type Present struct{
  ex int64
}

func (c *Present) Now() time.Time { return time.Now() }
func (c *Present) Duration() time.Duration { return time.Duration(c.ex) * time.Minute }


type Past struct{
  ex int64
}

func (c *Past) Now() time.Time {
  const longForm = "Jan 2, 2006 at 3:04pm (MST)"
  t, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")
  return t
}

func (c *Past) Duration() time.Duration { return time.Duration(c.ex) * time.Minute }


func TestSet(t *testing.T) {
  // when TTL is not defined
  key := "key"
  value := "value"

  moment := &Present{}
  Set(key, value, moment)

  expected := &cache.Item{"value", 0}

  cache := GetCache()
  actual := cache.Items[key]
  if *actual != *expected {
    t.Errorf("Set: expected %v, actual %v", expected, actual)
  }

  // when TTL is defined
  moment = &Present{30}

  Set(key, value, moment)

  actual = cache.Items[key]
  if actual.ExpireAt == 0 {
    t.Errorf("Set: ExpireAt should be greater than nil")
  }

}

func TestGet(t *testing.T) {
  key := "key"

  expected := "value"
  //setup
  Set(key, expected, &Present{})

  // run
  actual, _ := Get(key)

  if (actual != expected) {
    t.Errorf("Get: expected %v, actual %v", expected, actual)
  }
}



func TestUpd(t *testing.T) {
  key := "key"

  expected :="newValue"
  //setup
  Set(key, "value", &Present{})

  Upd(key, expected)
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
  Set(key, "value", &Present{})

  //run
  Del(key)

  _, found := Get(key)
  if (found != expected) {
    t.Errorf("Get: expected %v, actual %v", expected, found)
  }

}

func TestKeys(t *testing.T) {

  Set("adam[23]", "value", &Present{})

  pattern := "^[a-z]+[[0-9]+]$"

  keys := Keys(pattern)

  if (keys[0] != "adam[23]") {
    t.Errorf("Keys: expect %v, got %v", "adam[23]", keys[0])
  }

}

func TestDelRandomExpires(t *testing.T) {
  itemsToRemoveAmount := 5
  keys := []string{"key1", "key2", "key3"}
  for _, key := range keys {
    Set(key, "value", &Present{})
  }
  Set("key_to_expire", "value", &Past{10})

  DelRandomExpires(itemsToRemoveAmount)

  _, found := Get("key_to_expire")

  if (found) {
    t.Error("Item should be deleted")
  }

}

