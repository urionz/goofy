package contracts

import "time"

type Cache interface {
	Scan(key string, ptr interface{}, defVal ...interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	GetMultiple(keys []string, defVal interface{}) map[string]interface{}
	SetMultiple(values map[string]interface{}, ttl ...time.Duration) error
	DelMultiple(keys []string) error
	Has(key string) bool
}

type Store interface {
	Get(key string) (interface{}, error)
	Many(keys []string) []interface{}
	Put(key string, value interface{}, seconds time.Duration) error
	PutPure(key string, value interface{}, seconds time.Duration) error
	PutMany(kv map[string]interface{}, seconds int) error
	Increment(key string, value ...int64) error
	Decrement(key string, value ...int64) error
	Forever(key string, value interface{}) error
	ForeverPure(key string, value interface{}) error
	Forget(key string) error
	Has(key string) bool
	ItemKey(key string) string
	Flush() error
	GetPrefix() string
}

type CacheClosure = func() (interface{}, error)

type CacheRepository interface {
	Cache
	Tags(names ...string) (TaggableStore, error)
	Pull(key string, defVal ...interface{}) interface{}
	Put(key string, value interface{}, ttl time.Duration) error
	PutPure(key string, value interface{}, ttl time.Duration) error
	Add(key string, value interface{}, ttl ...time.Duration) error
	Increment(key string, value ...int64) error
	Decrement(key string, value ...int64) error
	Forever(key string, value interface{}) error
	ForeverPure(key string, value interface{}) error
	Remember(key string, ttl time.Duration, closure CacheClosure, ptr interface{}, force ...bool) error
	Sear(key string, closure CacheClosure, ptr interface{}, force ...bool) error
	RememberForever(key string, closure CacheClosure, ptr interface{}, force ...bool) error
	Forget(key string) error
	GetStore() Cache
}

type TaggableStore interface {
	CacheRepository
}

type CacheFactory interface {
	Store(name ...string) CacheRepository
}
