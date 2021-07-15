package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/urionz/goofy/contracts"
)

type Connection struct {
	client *contracts.RedisClient
	name   string
}

var _ contracts.RedisConnection = (*Connection)(nil)

func NewConnection(client *redis.Client) *Connection {
	return &Connection{
		client: client,
	}
}

func (conn *Connection) SetName(name string) *Connection {
	conn.name = name
	return conn
}

func (conn *Connection) GetName() string {
	return conn.name
}

func (conn *Connection) Client() *contracts.RedisClient {
	return conn.client
}

func (conn *Connection) Get(key string) string {
	return conn.client.Get(context.Background(), key).Val()
}

func (conn *Connection) Set(key string, value interface{}, expiration time.Duration) error {
	return conn.client.Set(context.Background(), key, value, expiration).Err()
}

func (conn *Connection) SetEX(key string, value interface{}, expiration time.Duration) error {
	return conn.client.SetEX(context.Background(), key, value, expiration).Err()
}

func (conn *Connection) Del(keys ...string) error {
	for _, key := range keys {
		if err := conn.client.Del(context.Background(), key).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (conn *Connection) Incr(key string) error {
	return conn.client.Incr(context.Background(), key).Err()
}

func (conn *Connection) Decr(key string) error {
	return conn.client.Decr(context.Background(), key).Err()
}

func (conn *Connection) IncrBy(key string, value int64) error {
	return conn.client.IncrBy(context.Background(), key, value).Err()
}

func (conn *Connection) DecrBy(key string, value int64) error {
	return conn.client.DecrBy(context.Background(), key, value).Err()
}

// 有序集合
func (conn *Connection) ZAdd(key string, members ...*contracts.Z) (int64, error) {
	return conn.client.ZAdd(context.Background(), key, members...).Result()
}

func (conn *Connection) ZCard(key string) (int64, error) {
	return conn.client.ZCard(context.Background(), key).Result()
}

func (conn *Connection) ZCount(key, min, max string) (int64, error) {
	return conn.client.ZCount(context.Background(), key, min, max).Result()
}

func (conn *Connection) ZIncrBy(key string, inc float64, member string) (float64, error) {
	return conn.client.ZIncrBy(context.Background(), key, inc, member).Result()
}

func (conn *Connection) ZInterStore(dest string, store *contracts.ZStore) (int64, error) {
	return conn.client.ZInterStore(context.Background(), dest, store).Result()
}

func (conn *Connection) ZLexCount(key, min, max string) (int64, error) {
	return conn.client.ZLexCount(context.Background(), key, min, max).Result()
}

func (conn *Connection) ZRange(key string, start, stop int64) ([]string, error) {
	return conn.client.ZRange(context.Background(), key, start, stop).Result()
}

func (conn *Connection) ZRangeByLex(key string, opt *contracts.ZRangeBy) ([]string, error) {
	return conn.client.ZRangeByLex(context.Background(), key, opt).Result()
}

func (conn *Connection) ZRangeByScore(key string, opt *contracts.ZRangeBy) ([]string, error) {
	return conn.client.ZRangeByScore(context.Background(), key, opt).Result()
}

func (conn *Connection) ZRank(key, member string) (int64, error) {
	return conn.client.ZRank(context.Background(), key, member).Result()
}

func (conn *Connection) ZRem(key, min, max string) (int64, error) {
	return conn.client.ZRem(context.Background(), key, min, max).Result()
}

func (conn *Connection) ZRemRangeByLex(key, min, max string) (int64, error) {
	return conn.client.ZRemRangeByLex(context.Background(), key, min, max).Result()
}

func (conn *Connection) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	return conn.client.ZRemRangeByRank(context.Background(), key, start, stop).Result()
}

func (conn *Connection) ZRemRangeByScore(key, min, max string) (int64, error) {
	return conn.client.ZRemRangeByScore(context.Background(), key, min, max).Result()
}

func (conn *Connection) ZRevRange(key string, start, stop int64) ([]string, error) {
	return conn.client.ZRevRange(context.Background(), key, start, stop).Result()
}

func (conn *Connection) ZRevRangeByScore(key string, opt *contracts.ZRangeBy) ([]string, error) {
	return conn.client.ZRevRangeByScore(context.Background(), key, opt).Result()
}

func (conn *Connection) ZRevRangeByLex(key string, opt *contracts.ZRangeBy) ([]string, error) {
	return conn.client.ZRevRangeByLex(context.Background(), key, opt).Result()
}

func (conn *Connection) ZRevRank(key, member string) (int64, error) {
	return conn.client.ZRevRank(context.Background(), key, member).Result()
}

func (conn *Connection) ZRevRangeByScoreWithScores(key string, opt *contracts.ZRangeBy) ([]contracts.Z, error) {
	return conn.client.ZRevRangeByScoreWithScores(context.Background(), key, opt).Result()
}

func (conn *Connection) ZRevRangeWithScores(key string, start, stop int64) ([]contracts.Z, error) {
	return conn.client.ZRevRangeWithScores(context.Background(), key, start, stop).Result()
}

func (conn *Connection) ZScore(key, member string) (float64, error) {
	return conn.client.ZScore(context.Background(), key, member).Result()
}

func (conn *Connection) ZScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return conn.client.ZScan(context.Background(), key, cursor, match, count).Result()
}

func (conn *Connection) ZIncr(key string, member *contracts.Z) (float64, error) {
	return conn.client.ZIncr(context.Background(), key, member).Result()
}

func (conn *Connection) ZIncrNX(key string, member *contracts.Z) (float64, error) {
	return conn.client.ZIncrNX(context.Background(), key, member).Result()
}

func (conn *Connection) ZIncrXX(key string, member *contracts.Z) (float64, error) {
	return conn.client.ZIncrXX(context.Background(), key, member).Result()
}

// 列表部分
func (conn *Connection) BLPop(timeout time.Duration, keys ...string) ([]string, error) {
	return conn.client.BLPop(context.Background(), timeout, keys...).Result()
}

func (conn *Connection) BRPop(timeout time.Duration, keys ...string) ([]string, error) {
	return conn.client.BRPop(context.Background(), timeout, keys...).Result()
}

func (conn *Connection) BRPopLPush(source, dest string, timeout time.Duration) (string, error) {
	return conn.client.BRPopLPush(context.Background(), source, dest, timeout).Result()
}

func (conn *Connection) LIndex(key string, index int64) (string, error) {
	return conn.client.LIndex(context.Background(), key, index).Result()
}

func (conn *Connection) LInsert(key, op string, pivot, value interface{}) (int64, error) {
	return conn.client.LInsert(context.Background(), key, op, pivot, value).Result()
}

func (conn *Connection) LLen(key string) (int64, error) {
	return conn.client.LLen(context.Background(), key).Result()
}

func (conn *Connection) LPop(key string) (string, error) {
	return conn.client.LPop(context.Background(), key).Result()
}

func (conn *Connection) LPush(key string, value ...interface{}) (int64, error) {
	return conn.client.LPush(context.Background(), key, value...).Result()
}

func (conn *Connection) LPushX(key string, value ...interface{}) (int64, error) {
	return conn.client.LPushX(context.Background(), key, value...).Result()
}

func (conn *Connection) LRange(key string, start, stop int64) ([]string, error) {
	return conn.client.LRange(context.Background(), key, start, stop).Result()
}

func (conn *Connection) LRem(key string, count int64, value interface{}) (int64, error) {
	return conn.client.LRem(context.Background(), key, count, value).Result()
}

func (conn *Connection) LSet(key string, index int64, value interface{}) (string, error) {
	return conn.client.LSet(context.Background(), key, index, value).Result()
}

func (conn *Connection) LTrim(key string, start, stop int64) (string, error) {
	return conn.client.LTrim(context.Background(), key, start, stop).Result()
}

func (conn *Connection) RPop(key string) (string, error) {
	return conn.client.RPop(context.Background(), key).Result()
}

func (conn *Connection) RPopLPush(source, dest string) (string, error) {
	return conn.client.RPopLPush(context.Background(), source, dest).Result()
}

func (conn *Connection) RPush(key string, value ...interface{}) (int64, error) {
	return conn.client.RPush(context.Background(), key, value...).Result()
}

func (conn *Connection) RPushX(key string, value ...interface{}) (int64, error) {
	return conn.client.RPushX(context.Background(), key, value...).Result()
}

// hash部分
func (conn *Connection) HSet(key string, value ...interface{}) error {
	return conn.client.HSet(context.Background(), key, value...).Err()
}

func (conn *Connection) HGet(key, field string) string {
	return conn.client.HGet(context.Background(), key, field).Val()
}

func (conn *Connection) HMSet(key string, value ...interface{}) error {
	return conn.client.HMSet(context.Background(), key, value...).Err()
}

func (conn *Connection) HMGet(key string, field ...string) []interface{} {
	return conn.client.HMGet(context.Background(), key, field...).Val()
}

func (conn *Connection) HDel(key string, fields ...string) (int64, error) {
	return conn.client.HDel(context.Background(), key, fields...).Result()
}

func (conn *Connection) HExists(key, field string) (bool, error) {
	return conn.client.HExists(context.Background(), key, field).Result()
}

func (conn *Connection) HGetAll(key string) (map[string]string, error) {
	return conn.client.HGetAll(context.Background(), key).Result()
}

func (conn *Connection) HIncrBy(key, field string, incr int64) (int64, error) {
	return conn.client.HIncrBy(context.Background(), key, field, incr).Result()
}

func (conn *Connection) HIncrByFloat(key, field string, incr float64) (float64, error) {
	return conn.client.HIncrByFloat(context.Background(), key, field, incr).Result()
}

func (conn *Connection) HKeys(key string) ([]string, error) {
	return conn.client.HKeys(context.Background(), key).Result()
}

func (conn *Connection) HLen(key string) (int64, error) {
	return conn.client.HLen(context.Background(), key).Result()
}

func (conn *Connection) HSetNX(key, field string, value interface{}) (bool, error) {
	return conn.client.HSetNX(context.Background(), key, field, value).Result()
}

func (conn *Connection) HVals(key string) ([]string, error) {
	return conn.client.HVals(context.Background(), key).Result()
}

func (conn *Connection) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return conn.client.HScan(context.Background(), key, cursor, match, count).Result()
}

// 集合
func (conn *Connection) SAdd(key string, members ...interface{}) error {
	return conn.client.SAdd(context.Background(), key, members...).Err()
}

func (conn *Connection) SCard(key string) (int64, error) {
	return conn.client.SCard(context.Background(), key).Result()
}

func (conn *Connection) SDiff(key ...string) ([]string, error) {
	return conn.client.SDiff(context.Background(), key...).Result()
}

func (conn *Connection) SDiffStore(dest string, keys ...string) (int64, error) {
	return conn.client.SDiffStore(context.Background(), dest, keys...).Result()
}

func (conn *Connection) SInter(key ...string) ([]string, error) {
	return conn.client.SInter(context.Background(), key...).Result()
}

func (conn *Connection) SInterStore(dest string, keys ...string) (int64, error) {
	return conn.client.SInterStore(context.Background(), dest, keys...).Result()
}

func (conn *Connection) SIsMember(key string, member interface{}) (bool, error) {
	return conn.client.SIsMember(context.Background(), key, member).Result()
}

func (conn *Connection) SMembers(key string) ([]string, error) {
	return conn.client.SMembers(context.Background(), key).Result()
}

func (conn *Connection) SMove(source, dest string, member interface{}) (bool, error) {
	return conn.client.SMove(context.Background(), source, dest, member).Result()
}

func (conn *Connection) SPop(key string) (string, error) {
	return conn.client.SPop(context.Background(), key).Result()
}

func (conn *Connection) SRandMember(key string) (string, error) {
	return conn.client.SRandMember(context.Background(), key).Result()
}

func (conn *Connection) SRandMemberN(key string, count int64) ([]string, error) {
	return conn.client.SRandMemberN(context.Background(), key, count).Result()
}

func (conn *Connection) SRem(key string, members ...interface{}) error {
	return conn.client.SRem(context.Background(), key, members...).Err()
}

func (conn *Connection) SUnion(keys ...string) ([]string, error) {
	return conn.client.SUnion(context.Background(), keys...).Result()
}

func (conn *Connection) SUnionStore(dest string, keys ...string) (int64, error) {
	return conn.client.SUnionStore(context.Background(), dest, keys...).Result()
}

func (conn *Connection) SScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return conn.client.SScan(context.Background(), key, cursor, match, count).Result()
}

func (conn *Connection) Multi(cb contracts.MultiFunc) error {
	tx := conn.client.TxPipeline()
	cb(tx)
	if _, err := tx.Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
