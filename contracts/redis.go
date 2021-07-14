package contracts

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisFactory interface {
	Connection(name ...string) (RedisConnection, error)
}

type (
	RedisClient = redis.Client
	Pipeliner   = redis.Pipeliner
	Pipeline    = redis.Pipeline
	MultiFunc   func(pipe Pipeliner)
	Z           = redis.Z
	ZRangeBy    = redis.ZRangeBy
)

type RedisConnection interface {
	GetName() string
	Client() *RedisClient
	Get(key string) string
	Set(key string, value interface{}, expiration time.Duration) error
	Del(keys ...string) error
	SetEX(key string, value interface{}, expiration time.Duration) error
	ZAdd(key string, members ...*Z) (int64, error)
	ZScore(key, member string) (float64, error)
	ZLexCount(key, min, max string) (int64, error)
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
	ZRange(key string, start, stop int64) ([]string, error)
	ZRangeByLex(key string, opt *ZRangeBy) ([]string, error)
	ZRangeByScore(key string, opt *ZRangeBy) ([]string, error)
	ZCount(key, min, max string) (int64, error)
	ZRem(key, min, max string) (int64, error)
	ZRemRangeByLex(key, min, max string) (int64, error)
	ZRemRangeByRank(key string, start, stop int64) (int64, error)
	ZRemRangeByScore(key, min, max string) (int64, error)
	ZRevRange(key string, start, stop int64) ([]string, error)
	ZRevRangeByLex(key string, opt *ZRangeBy) ([]string, error)
	ZRevRangeByScore(key string, opt *ZRangeBy) ([]string, error)
	ZRevRangeByScoreWithScores(key string, opt *ZRangeBy) ([]Z, error)
	ZRevRangeWithScores(key string, start, stop int64) ([]Z, error)
	ZRevRank(key, member string) (int64, error)
	ZScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error)
	ZRank(key, member string) (int64, error)
	ZCard(key string) (int64, error)
	ZIncr(key string, member *Z) (float64, error)
	ZIncrBy(key string, inc float64, member string) (float64, error)
	ZIncrNX(key string, member *Z) (float64, error)
	ZIncrXX(key string, member *Z) (float64, error)
	Multi(cb MultiFunc) error
}
