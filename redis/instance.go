package redis

import (
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
