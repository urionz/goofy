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
