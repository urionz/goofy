package contracts

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisFactory interface {
	Connection(name ...string) (RedisConnection, error)
}

type MultiFunc func(pipe redis.Pipeliner)

type RedisConnection interface {
	GetName() string
	Get(key string) string
	Set(key string, value interface{}, expiration time.Duration) error
	Del(keys ...string) error
	SetEX(key string, value interface{}, expiration time.Duration) error
	SAdd(key string, members ...interface{}) error
	Multi(cb MultiFunc) error
}
