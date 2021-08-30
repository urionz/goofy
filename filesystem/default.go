package filesystem

import "github.com/urionz/goofy/contracts"

var instance *Manager

func Disk(name ...string) contracts.Filesystem {
	return instance.Disk(name...)
}
