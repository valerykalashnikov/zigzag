package zigzag

import (
	"github.com/valerykalashnikov/zigzag/cache"
)

type Importer interface {
	Import() (map[string]*cache.Item, error)
}

func New(cacheType, is_slave string) (db *DB, err error) {
	var (
		store cache.DataStore
		slave bool
	)
	store, err = cache.CreateCache(cacheType)
	switch is_slave {
	case "1":
		slave = true
	default:
		slave = false
	}
	db = &DB{store, slave}
	return
}

type DB struct {
	store    cache.DataStore
	is_slave bool
}

func (db *DB) Set(key string, value interface{}, m cache.Momenter) {
	db.store.Set(key, value, m)
}

func (db *DB) Get(key string) (interface{}, bool) {
	if item, found := db.store.Get(key); found {
		if item.Expired() {
			db.store.Del(key)
			return nil, false
		}
		return item.Object, true
	}
	return nil, false
}

func (db *DB) Upd(key string, value interface{}) bool {
	return db.store.Upd(key, value)
}

func (db *DB) Del(key string) {
	db.store.Del(key)
}

func (db *DB) Keys(pattern string) []string {
	return db.store.Keys(pattern)
}

func (db *DB) Items() map[string]*cache.Item {
	return db.store.Items()
}

func (db *DB) DelRandomExpires(num int) int {
	items := db.store.Items()
	length := len(items)
	expiresRemoved := 0
	i := 0
	// The Go runtime actually randomizes the map iteration order
	for k, v := range items {
		if i == length {
			return expiresRemoved
		}
		if i == num {
			correlation := float64(expiresRemoved) / float64(i)
			if correlation < 0.25 {
				return expiresRemoved
			} else {
				num *= 2
			}
		}
		if v.Expired() {
			db.store.Del(k)
			expiresRemoved += 1
		}
		i += 1
	}
	return expiresRemoved
}
