package zigzag

import (
	"github.com/valerykalashnikov/zigzag/cache"
	"testing"
	"time"
	// "fmt"
)

type Present struct {
	ex int64
}

func (c *Present) Now() time.Time          { return time.Now() }
func (c *Present) Duration() time.Duration { return time.Duration(c.ex) * time.Minute }

type Past struct {
	ex int64
}

func (c *Past) Now() time.Time {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	t, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")
	return t
}

func (c *Past) Duration() time.Duration { return time.Duration(c.ex) * time.Minute }

func ConnectDB(engType, slavery string, t *testing.T) (db *DB) {
	db, err := New(engType, slavery)
	if err != nil {
		t.Errorf("Set: expected nil, got error, %v", err)
	}
	return
}

func Check(status, expected bool, t *testing.T) error {
	actual := status
	if actual != expected {
		t.Errorf("Set: expected %v, actual %v", expected, actual)
	}
	return nil
}

func TestNew(t *testing.T) {
	// test master
	db := ConnectDB("cache", "0", t)
	_ = Check(db.is_slave, false, t)

	// test slave
	db = ConnectDB("cache", "1", t)
	_ = Check(db.is_slave, true, t)
}

func TestCheckSlavery(t *testing.T) {
	// test master
	db := ConnectDB("cache", "0", t)
	_ = Check(db.CheckSlavery(), false, t)

	// test slave
	db = ConnectDB("cache", "1", t)
	_ = Check(db.CheckSlavery(), true, t)
}

func TestSet(t *testing.T) {
	var moment cache.Momenter
	// when TTL is not defined
	db, err := New("sharded", "0")
	if err != nil {
		t.Errorf("Set: expected nil, got error, %v", err)
	}
	key := "key"
	value := "value"

	moment = &Present{}
	db.Set(key, value, moment)

	expected := "value"

	actual, _ := db.Get(key)
	if actual != expected {
		t.Errorf("Set: expected %v, actual %v", expected, actual)
	}

	// when TTL is defined
	moment = &Present{30}

	db.Set(key, value, moment)

	actual, _ = db.Get(key)
	if actual == nil {
		t.Errorf("Set: item shouldn't be expired")
	}
}

func TestGet(t *testing.T) {
	var moment cache.Momenter
	db, err := New("sharded", "0")

	if err != nil {
		t.Errorf("Set: expected nil, got error, %v", err)
	}

	key := "key1"

	expected := "value"
	moment = &Present{}

	db.Set(key, expected, moment)

	actual, _ := db.Get(key)

	if actual != expected {
		t.Errorf("Get: expected %v, actual %v", expected, actual)
	}

	// when item is expired
	moment = &Past{10}

	db.Set(key, expected, moment)

	actual, _ = db.Get(key)

	if actual != nil {
		t.Errorf("Get: item should be expired")
	}
}

func TestUpd(t *testing.T) {
	db, _ := New("sharded", "0")

	key := "key"

	expected := "newValue"
	//setup
	db.Set(key, "value", &Present{})

	db.Upd(key, expected)
	// run
	actual, _ := db.Get(key)

	if actual != expected {
		t.Errorf("Get: expected %v, actual %v", expected, actual)
	}
}

func TestDel(t *testing.T) {
	db, _ := New("sharded", "0")

	key := "key"

	expected := false
	//setup
	db.Set(key, "value", &Present{})

	//run
	db.Del(key)

	_, found := db.Get(key)
	if found != expected {
		t.Errorf("Get: expected %v, actual %v", expected, found)
	}

}

func TestKeys(t *testing.T) {
	db, _ := New("sharded", "0")

	db.Set("adam[23]", "value", &Present{})

	pattern := "^[a-z]+[[0-9]+]$"

	keys := db.Keys(pattern)

	if keys[0] != "adam[23]" {
		t.Errorf("Keys: expect %v, got %v", "adam[23]", keys[0])
	}

}

func TestDelRandomExpires(t *testing.T) {
	db, _ := New("sharded", "0")

	itemsToRemoveAmount := 5
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		db.Set(key, "value", &Present{})
	}
	db.Set("key_to_expire", "value", &Past{10})

	db.DelRandomExpires(itemsToRemoveAmount)

	_, found := db.Get("key_to_expire")

	if found {
		t.Error("Item should be deleted")
	}

}
