package cache

import (
	"time"

	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goutil/jsonutil"
)

type RedisStore struct {
	redis      contracts.RedisFactory
	prefix     string
	connection string
	TaggableStore
}

func NewRedisStore(redis contracts.RedisFactory, prefix, connection string) *RedisStore {
	return &RedisStore{
		redis:      redis,
		prefix:     prefix,
		connection: connection,
	}
}

func (r *RedisStore) Get(key string) (interface{}, error) {
	var standerValue StanderValue
	value := r.Connection().Get(r.prefix + key)
	if value == "" {
		return nil, nil
	}
	if err := jsonutil.Decode([]byte(value), &standerValue); err != nil {
		return value, err
	}
	return standerValue.Value, nil
}

func (r *RedisStore) Set(key string, value interface{}, ttl time.Duration) error {
	return r.Put(key, value, ttl)
}

func (r *RedisStore) Put(key string, value interface{}, seconds time.Duration) error {
	var err error
	var raw []byte
	if raw, err = jsonutil.Encode(StanderValue{Value: value}); err != nil {
		return err
	}
	return r.Connection().Set(key, string(raw), seconds)
}

func (r *RedisStore) PutInt(key string, value int64, seconds time.Duration) error {
	return r.Connection().Set(key, value, seconds)
}

func (r *RedisStore) Forever(key string, value interface{}) error {
	return r.Set(key, value, 0)
}

func (r *RedisStore) ForeverInt(key string, value int64) error {
	return r.PutInt(key, value, 0)
}

func (r *RedisStore) Forget(key string) error {
	return r.Connection().Del(key)
}

func (r *RedisStore) Has(key string) bool {
	return r.Connection().Exists(key)
}

func (r *RedisStore) Tags(names ...string) (contracts.TaggableStore, error) {
	return NewRedisTaggedCache(r, NewTagSet(r, names...)), nil
}

func (r *RedisStore) Increment(key string, steps ...int64) error {
	var step int64 = 1
	if len(steps) == 0 {
		step = steps[0]
	}
	return r.Connection().IncrBy(key, step)
}
func (r *RedisStore) Decrement(key string, steps ...int64) error {
	var step int64 = 1
	if len(steps) == 0 {
		step = steps[0]
	}
	return r.Connection().DecrBy(key, step)
}

func (r *RedisStore) ItemKey(key string) string {
	return key
}

func (r *RedisStore) Connection() contracts.RedisConnection {
	rds, err := r.redis.Connection(r.connection)
	if err != nil {
		panic(err)
	}
	return rds
}
