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

func Multi(cb contracts.MultiFunc, name ...string) error {
	return Conn(name...).Multi(cb)
}
