package contracts

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisFactory interface {
	Connection(name ...string) (RedisConnection, error)
}

type RedisClient = redis.Client

type MultiFunc func(pipe redis.Pipeliner)

type RedisConnection interface {
	GetName() string
	Client() *RedisClient
	Get(key string) string
	Set(key string, value interface{}, expiration time.Duration) error
	Del(keys ...string) error
	SetEX(key string, value interface{}, expiration time.Duration) error
	SAdd(key string, members ...interface{}) error
	Incr(key string) error
	Decr(key string) error
	IncrBy(key string, value int64) error
	DecrBy(key string, value int64) error
	HSet(key string, value ...interface{}) error
	HGet(key, field string) string
	HMSet(key string, value ...interface{}) error
	HMGet(key string, field ...string) []interface{}
	SIsMember(key string, member interface{}) (bool, error)
	SRem(key string, members ...interface{}) error
	Multi(cb MultiFunc) error
}
