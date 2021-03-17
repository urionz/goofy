package cache

import "github.com/urionz/goofy/contracts"

var instance *Manager

func Store(name ...string) contracts.CacheRepository {
	return instance.Store(name...)
}
