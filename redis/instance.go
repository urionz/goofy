package redis

import (
	"time"

	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/log"
)

var instance *Manager

func Conn(name ...string) contracts.RedisConnection {
	conn, err := instance.Connection()
	if err != nil {
		log.Panic(err)
	}
	return conn
}

func Multi(cb contracts.MultiFunc, conn ...string) error {
	return Conn(conn...).Multi(cb)
}

func Get(key string, conn ...string) string {
	return Conn(conn...).Get(key)
}

func Set(key string, value interface{}, expiration time.Duration, conn ...string) error {
	return Conn(conn...).Set(key, value, expiration)
}

func SetEX(key string, value interface{}, expiration time.Duration, conn ...string) error {
	return Conn(conn...).SetEX(key, value, expiration)
}

func Incr(key string, conn ...string) error {
	return Conn(conn...).Incr(key)
}

func Decr(key string, conn ...string) error {
	return Conn(conn...).Decr(key)
}

func IncrBy(key string, value int64, conn ...string) error {
	return Conn(conn...).IncrBy(key, value)
}

func DecrBy(key string, value int64, conn ...string) error {
	return Conn(conn...).DecrBy(key, value)
}

func HGet(key, field string, conn ...string) string {
	return Conn(conn...).HGet(key, field)
}

func SIsMember(key string, member interface{}, conn ...string) (bool, error) {
	return Conn(conn...).SIsMember(key, member)
}

func ZAdd(key string, members []*contracts.Z, conn ...string) (int64, error) {
	return Conn(conn...).ZAdd(key, members...)
}

func ZScore(key, member string, conn ...string) (float64, error) {
	return Conn(conn...).ZScore(key, member)
}

func ZLexCount(key, min, max string, conn ...string) (int64, error) {
	return Conn(conn...).ZLexCount(key, min, max)
}

func ZRange(key string, start, stop int64, conn ...string) ([]string, error) {
	return Conn(conn...).ZRange(key, start, stop)
}

func ZRangeByLex(key string, opt *contracts.ZRangeBy, conn ...string) ([]string, error) {
	return Conn(conn...).ZRangeByLex(key, opt)
}

func ZRangeByScore(key string, opt *contracts.ZRangeBy, conn ...string) ([]string, error) {
	return Conn(conn...).ZRangeByScore(key, opt)
}

func ZCount(key, min, max string, conn ...string) (int64, error) {
	return Conn(conn...).ZCount(key, min, max)
}

func ZRem(key, min, max string, conn ...string) (int64, error) {
	return Conn(conn...).ZRem(key, min, max)
}

func ZRemRangeByLex(key, min, max string, conn ...string) (int64, error) {
	return Conn(conn...).ZRemRangeByLex(key, min, max)
}

func ZRemRangeByRank(key string, start, stop int64, conn ...string) (int64, error) {
	return Conn(conn...).ZRemRangeByRank(key, start, stop)
}

func ZRemRangeByScore(key, min, max string, conn ...string) (int64, error) {
	return Conn(conn...).ZRemRangeByScore(key, min, max)
}

func ZRevRange(key string, start, stop int64, conn ...string) ([]string, error) {
	return Conn(conn...).ZRevRange(key, start, stop)
}

func ZRevRangeByLex(key string, opt *contracts.ZRangeBy, conn ...string) ([]string, error) {
	return Conn(conn...).ZRevRangeByLex(key, opt)
}

func ZRevRangeByScore(key string, opt *contracts.ZRangeBy, conn ...string) ([]string, error) {
	return Conn(conn...).ZRevRangeByScore(key, opt)
}

func ZRevRangeByScoreWithScores(key string, opt *contracts.ZRangeBy, conn ...string) ([]contracts.Z, error) {
	return Conn(conn...).ZRevRangeByScoreWithScores(key, opt)
}

func ZRevRangeWithScores(key string, start, stop int64, conn ...string) ([]contracts.Z, error) {
	return Conn(conn...).ZRevRangeWithScores(key, start, stop)
}

func ZRevRank(key, member string, conn ...string) (int64, error) {
	return Conn(conn...).ZRevRank(key, member)
}

func ZScan(key string, cursor uint64, match string, count int64, conn ...string) ([]string, uint64, error) {
	return Conn(conn...).ZScan(key, cursor, match, count)
}

func ZRank(key, member string, conn ...string) (int64, error) {
	return Conn(conn...).ZRank(key, member)
}

func ZCard(key string, conn ...string) (int64, error) {
	return Conn(conn...).ZCard(key)
}

func ZIncr(key string, member *contracts.Z, conn ...string) (float64, error) {
	return Conn(conn...).ZIncr(key, member)
}

func ZIncrBy(key string, inc float64, member string, conn ...string) (float64, error) {
	return Conn(conn...).ZIncrBy(key, inc, member)
}

func ZIncrNX(key string, member *contracts.Z, conn ...string) (float64, error) {
	return Conn(conn...).ZIncrNX(key, member)
}

func ZIncrXX(key string, member *contracts.Z, conn ...string) (float64, error) {
	return Conn(conn...).ZIncrXX(key, member)
}

func BLPop(timeout time.Duration, keys []string, conn ...string) ([]string, error) {
	return Conn(conn...).BLPop(timeout, keys...)
}

func BRPop(timeout time.Duration, keys []string, conn ...string) ([]string, error) {
	return Conn(conn...).BRPop(timeout, keys...)
}

func BRPopLPush(source, dest string, timeout time.Duration, conn ...string) (string, error) {
	return Conn(conn...).BRPopLPush(source, dest, timeout)
}

func LIndex(key string, index int64, conn ...string) (string, error) {
	return Conn(conn...).LIndex(key, index)
}

func LInsert(key, op string, pivot, value interface{}, conn ...string) (int64, error) {
	return Conn(conn...).LInsert(key, op, pivot, value)
}

func LLen(key string, conn ...string) (int64, error) {
	return Conn(conn...).LLen(key)
}

func LPop(key string, conn ...string) (string, error) {
	return Conn(conn...).LPop(key)
}

func LPush(key string, value []interface{}, conn ...string) (int64, error) {
	return Conn(conn...).LPush(key, value)
}

func LPushX(key string, value []interface{}, conn ...string) (int64, error) {
	return Conn(conn...).LPushX(key, value)
}

func LRange(key string, start, stop int64, conn ...string) ([]string, error) {
	return Conn(conn...).LRange(key, start, stop)
}

func LRem(key string, count int64, value interface{}, conn ...string) (int64, error) {
	return Conn(conn...).LRem(key, count, value)
}

func LSet(key string, index int64, value interface{}, conn ...string) (string, error) {
	return Conn(conn...).LSet(key, index, value)
}

func LTrim(key string, start, stop int64, conn ...string) (string, error) {
	return Conn(conn...).LTrim(key, start, stop)
}

func RPop(key string, conn ...string) (string, error) {
	return Conn(conn...).RPop(key)
}

func RPopLPush(source, dest string, conn ...string) (string, error) {
	return Conn(conn...).RPopLPush(source, dest)
}

func RPush(key string, value []interface{}, conn ...string) (int64, error) {
	return Conn(conn...).RPush(key, value...)
}

func RPushX(key string, value []interface{}, conn ...string) (int64, error) {
	return Conn(conn...).RPushX(key, value...)
}

func HMGet(key string, fields []string, conn ...string) []interface{} {
	return Conn(conn...).HMGet(key, fields...)
}

func HMSet(key string, values []interface{}, conn ...string) error {
	return Conn(conn...).HMSet(key, values...)
}

func HSet(key string, values []interface{}, conn ...string) error {
	return Conn(conn...).HSet(key, values...)
}

func HDel(key string, fields []string, conn ...string) (int64, error) {
	return Conn(conn...).HDel(key, fields...)
}

func HExists(key, field string, conn ...string) (bool, error) {
	return Conn(conn...).HExists(key, field)
}

func HGetAll(key string, conn ...string) (map[string]string, error) {
	return Conn(conn...).HGetAll(key)
}

func HIncrBy(key, field string, incr int64, conn ...string) (int64, error) {
	return Conn(conn...).HIncrBy(key, field, incr)
}

func HIncrByFloat(key, field string, incr float64, conn ...string) (float64, error) {
	return Conn(conn...).HIncrByFloat(key, field, incr)
}

func HKeys(key string, conn ...string) ([]string, error) {
	return Conn(conn...).HKeys(key)
}

func HLen(key string, conn ...string) (int64, error) {
	return Conn(conn...).HLen(key)
}

func HSetNX(key, field string, value interface{}, conn ...string) (bool, error) {
	return Conn(conn...).HSetNX(key, field, value)
}

func HVals(key string, conn ...string) ([]string, error) {
	return Conn(conn...).HVals(key)
}

func HScan(key string, cursor uint64, match string, count int64, conn ...string) ([]string, uint64, error) {
	return Conn(conn...).HScan(key, cursor, match, count)
}
