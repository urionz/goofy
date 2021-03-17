package contracts

import "time"

type RedisFactory interface {
	Connection(name ...string) (RedisConnection, error)
}

type RedisConnection interface {
	GetName() string
	Get(key string) string
	Set(key string, value interface{}, expiration time.Duration) error
	SetEX(key string, value interface{}, expiration time.Duration) error
	SAdd(key string, members ...interface{}) error
}
