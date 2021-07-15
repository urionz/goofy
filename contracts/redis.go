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
	ZStore      = redis.ZStore
)

type RedisConnection interface {
	GetName() string
	Client() *RedisClient
	Get(key string) string
	Set(key string, value interface{}, expiration time.Duration) error
	Del(keys ...string) error
	SetEX(key string, value interface{}, expiration time.Duration) error

	ZAdd(key string, members ...*Z) (int64, error)
	ZCard(key string) (int64, error)
	ZCount(key, min, max string) (int64, error)
	ZIncrBy(key string, inc float64, member string) (float64, error)
	ZInterStore(dest string, store *ZStore) (int64, error)

	ZScore(key, member string) (float64, error)
	ZLexCount(key, min, max string) (int64, error)
	SAdd(key string, members ...interface{}) error
	Incr(key string) error
	Decr(key string) error
	IncrBy(key string, value int64) error
	DecrBy(key string, value int64) error
	SIsMember(key string, member interface{}) (bool, error)
	SRem(key string, members ...interface{}) error
	ZRange(key string, start, stop int64) ([]string, error)
	ZRangeByLex(key string, opt *ZRangeBy) ([]string, error)
	ZRangeByScore(key string, opt *ZRangeBy) ([]string, error)

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
	ZIncr(key string, member *Z) (float64, error)
	ZIncrNX(key string, member *Z) (float64, error)
	ZIncrXX(key string, member *Z) (float64, error)
	Multi(cb MultiFunc) error

	// hash
	HGet(key, field string) string
	HMGet(key string, field ...string) []interface{}
	HMSet(key string, value ...interface{}) error
	HSet(key string, value ...interface{}) error
	HDel(key string, fields ...string) (int64, error)
	HExists(key, field string) (bool, error)
	HGetAll(key string) (map[string]string, error)
	HIncrBy(key, field string, incr int64) (int64, error)
	HIncrByFloat(key, field string, incr float64) (float64, error)
	HKeys(key string) ([]string, error)
	HLen(key string) (int64, error)
	HSetNX(key, field string, value interface{}) (bool, error)
	HVals(key string) ([]string, error)
	HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error)

	// 列表
	BLPop(timeout time.Duration, keys ...string) ([]string, error)
	BRPop(timeout time.Duration, keys ...string) ([]string, error)
	BRPopLPush(source, dest string, timeout time.Duration) (string, error)
	LIndex(key string, index int64) (string, error)
	LInsert(key, op string, pivot, value interface{}) (int64, error)
	LLen(key string) (int64, error)
	LPop(key string) (string, error)
	LPush(key string, value ...interface{}) (int64, error)
	LPushX(key string, value ...interface{}) (int64, error)
	LRange(key string, start, stop int64) ([]string, error)
	LRem(key string, count int64, value interface{}) (int64, error)
	LSet(key string, index int64, value interface{}) (string, error)
	LTrim(key string, start, stop int64) (string, error)
	RPop(key string) (string, error)
	RPopLPush(source, dest string) (string, error)
	RPush(key string, value ...interface{}) (int64, error)
	RPushX(key string, value ...interface{}) (int64, error)
}
