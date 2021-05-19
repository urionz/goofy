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

func (conn *Connection) SAdd(key string, members ...interface{}) error {
	return conn.client.SAdd(context.Background(), key, members...).Err()
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

func (conn *Connection) SIsMember(key string, member interface{}) (bool, error) {
	return conn.client.SIsMember(context.Background(), key, member).Result()
}

func (conn *Connection) SRem(key string, members ...interface{}) error {
	return conn.client.SRem(context.Background(), key, members...).Err()
}

func (conn *Connection) Multi(cb contracts.MultiFunc) error {
	tx := conn.client.TxPipeline()
	cb(tx)
	if _, err := tx.Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
