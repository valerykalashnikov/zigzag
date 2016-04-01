package cache

import (
        "reflect"
        "time"
        "testing"
        )
var cache *Cache

type AncientTime struct {
  duration time.Duration
}

func (c *AncientTime) Now() time.Time {
  const longForm = "Jan 2, 2006 at 3:04pm (MST)"
  t, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")
  return t
}
func (c *AncientTime) Duration() time.Duration { return c.duration }


type PresentTime struct {
  duration time.Duration
}

func (c *PresentTime) Now() time.Time {
  return time.Now()
}
func (c *PresentTime) Duration() time.Duration { return c.duration }

var ancientTime *AncientTime

var presentTime *PresentTime

func setup() {
  cache = &Cache{
    Items: make(map[string]*Item),
  }
}

func TestSet(t *testing.T) {

  setup()
  // when TTL is not defined
  presentTime = &PresentTime{0}
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
    cache.Set(tt.key, tt.value, presentTime)
    actual := cache.Items[tt.key]
    if !reflect.DeepEqual(actual, tt.expected) {
      t.Errorf("Set: expected %v, actual %v", tt.expected, actual)
    }
  }

  //when TTL is defined
  ex := int64(30)
  key :=   "key"
  value := "value"
  duration := time.Duration(ex) * time.Minute
  ancientTime = &AncientTime{duration}
  cache.Set(key, value, ancientTime)
  actual := cache.Items[key]
  if actual.ExpireAt == 0 {
    t.Errorf("Set: ExpireAt should be greater than nil")
  }
  if actual.ExpireAt != ancientTime.Now().Add(duration).UnixNano() {
    t.Errorf("Set: ExpireAt should be now + duration")
  }

}

func TestGet(t *testing.T) {
  setup()
  // when ttl is not defined
  presentTime = &PresentTime{0}
  key := "key"
  expected := &Item{"value", 0}

  cache.Set(key, "value", presentTime)
  actual, _ := cache.Get(key)

  if (*actual != *expected) {
    t.Errorf("Get: expected %v, actual %v", expected, actual)
  }

  //when key is not present in storage
  delete(cache.Items, "key")

  actual, found := cache.Get(key)
  if (found != false) {
    t.Error("Get: expected false while returning empty value, found flag is", found)
  }

  //when key is presented in storage but it is outdated
  ancientTime = &AncientTime{10}
  delete(cache.Items, "key")
  cache.Set(key, "value", ancientTime)
  actual, found = cache.Get(key)
  if (found == true) {
    t.Error("Get: Expired value shouldn't be returned, actual =", found)
  }
}

func TestUpd(t *testing.T) {
  setup()

  //it shouldn't change expiration
  key :=   "key"
  value := "value"
  ex := int64(30)

  duration := time.Duration(ex) * time.Minute
  ancientTime = &AncientTime{duration}

  cache.Set(key, value, ancientTime)
  expectedValue := "newValue"
  cache.Upd(key, expectedValue)
  actualValue := cache.Items[key].Object

  if (actualValue != expectedValue) {
    t.Errorf("Upd: Object should be changed, expect %v, actual %v", expectedValue, actualValue)
  }

}

func TestDel(t *testing.T) {
  setup()

  key := "key"

  expected := false

  presentTime = &PresentTime{0}
  cache.Set(key, "value", presentTime)

  cache.Del(key)

  _, found := cache.Get(key)

  if (found != expected) {
    t.Errorf("Get: expected %v, actual %v", expected, found)
  }

}

func TestKeys(t *testing.T) {
  setup()

  presentTime = &PresentTime{0}
  cache.Set("adam[23]", "value", presentTime)
  cache.Set("eve[7]", "value", presentTime)
  cache.Set("Job[48]", "value", presentTime)
  cache.Set("snakey", "value", presentTime)

  pattern := "^[a-z]+[[0-9]+]$"

  keys := cache.Keys(pattern)

  if (!stringInSlice("adam[23]", keys)) {
    t.Errorf("Keys: expect %v, got %v", "adam[23]", keys[0])
  }
  if (!stringInSlice("eve[7]", keys)) {
    t.Errorf("Keys: expect %v, got %v", "eve[7]", keys[1])
  }

  if (len(keys) != 2) {
    t.Error("Keys: length of keys should be 2")
  }
}


func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}
