package cache

import (
	"fmt"
	"regexp"
	"sync"
)

// T table for Pearson hashing from RFC 3074.
var T = [256]uint8{
	251, 175, 119, 215, 81, 14, 79, 191, 103, 49, 181, 143, 186, 157, 0,
	232, 31, 32, 55, 60, 152, 58, 17, 237, 174, 70, 160, 144, 220, 90, 57,
	223, 59, 3, 18, 140, 111, 166, 203, 196, 134, 243, 124, 95, 222, 179,
	197, 65, 180, 48, 36, 15, 107, 46, 233, 130, 165, 30, 123, 161, 209, 23,
	97, 16, 40, 91, 219, 61, 100, 10, 210, 109, 250, 127, 22, 138, 29, 108,
	244, 67, 207, 9, 178, 204, 74, 98, 126, 249, 167, 116, 34, 77, 193,
	200, 121, 5, 20, 113, 71, 35, 128, 13, 182, 94, 25, 226, 227, 199, 75,
	27, 41, 245, 230, 224, 43, 225, 177, 26, 155, 150, 212, 142, 218, 115,
	241, 73, 88, 105, 39, 114, 62, 255, 192, 201, 145, 214, 168, 158, 221,
	148, 154, 122, 12, 84, 82, 163, 44, 139, 228, 236, 205, 242, 217, 11,
	187, 146, 159, 64, 86, 239, 195, 42, 106, 198, 118, 112, 184, 172, 87,
	2, 173, 117, 176, 229, 247, 253, 137, 185, 99, 164, 102, 147, 45, 66,
	231, 52, 141, 211, 194, 206, 246, 238, 56, 110, 78, 248, 63, 240, 189,
	93, 92, 51, 53, 183, 19, 171, 72, 50, 33, 104, 101, 69, 8, 252, 83, 120,
	76, 135, 85, 54, 202, 125, 188, 213, 96, 235, 136, 208, 162, 129, 190,
	132, 156, 38, 47, 1, 7, 254, 24, 4, 216, 131, 89, 21, 28, 133, 37, 153,
	149, 80, 170, 68, 6, 169, 234, 151,
}

func phash(key string) (h byte) {
	for _, c := range []byte(key) {
		h = T[h^c]
	}
	return
}

type ShardedCache map[string]*CacheShard

type CacheShard struct {
	items map[string]*Item
	mux   *sync.RWMutex
}

func (c ShardedCache) Set(key string, value interface{}, moment Momenter) {
	var expireAt int64
	duration := moment.Duration()
	if duration != 0 {
		expireAt = moment.Now().Add(duration).UnixNano()
	}
	shard := c.getShard(key)
	shard.mux.Lock()
	defer shard.mux.Unlock()
	shard.items[key] = &Item{
		Object:   value,
		ExpireAt: expireAt,
	}
}

func (c ShardedCache) Get(key string) (*Item, bool) {
	shard := c.getShard(key)
	shard.mux.RLock()
	defer shard.mux.RUnlock()

	if item, ok := shard.items[key]; ok {
		return item, true
	}
	return nil, false
}

func (c ShardedCache) Upd(key string, newValue interface{}) bool {
	// it saves TTL
	shard := c.getShard(key)
	shard.mux.Lock()
	defer shard.mux.Unlock()
	item, found := shard.items[key]
	if !found {
		return false
	}
	item.Object = newValue
	return true
}

func (c ShardedCache) Del(key string) {
	shard := c.getShard(key)
	shard.mux.Lock()
	defer shard.mux.Unlock()
	delete(shard.items, key)
}

func (c ShardedCache) Keys(pattern string) []string {
	var keys []string
	var validKey = regexp.MustCompile(pattern)
	for _, shard := range c {

		for k := range shard.items {
			if validKey.MatchString(k) {
				keys = append(keys, k)
			}
		}
	}

	return keys
}

func (c ShardedCache) Items() map[string]*Item {
	results := make(map[string]*Item)
	for _, shard := range c {
		for key, item := range shard.items {
			results[key] = item
		}
	}
	return results
}

func (c ShardedCache) getShard(key string) *CacheShard {
	shardKey := fmt.Sprintf("%02x", phash(key))
	return c[shardKey]
}

func NewShardedCache() DataStore {
	c := make(ShardedCache, 256)
	for i := 0; i < 256; i++ {
		c[fmt.Sprintf("%02x", i)] = &CacheShard{
			items: make(map[string]*Item, 2048),
			mux:   new(sync.RWMutex),
		}
	}
	return c
}
