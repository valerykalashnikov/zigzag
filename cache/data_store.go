package cache

type DataStore interface {
	Set(key string, value interface{}, moment Momenter)
	Get(key string) (*Item, bool)
	Upd(key string, newValue interface{}) bool
	Del(key string)
	Keys(pattern string) []string
	Items() map[string]*Item
}
