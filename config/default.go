package config

import (
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goutil"
)

var serve *Configure

func LoadExists(files ...string) error {
	return serve.LoadExists(files...)
}

func Data() map[string]interface{} {
	return serve.Data()
}

func Get(key string, findByPath ...bool) interface{} {
	return serve.Get(key, findByPath...)
}

func Exists(key string, findByPath ...bool) bool {
	return serve.Exists(key, findByPath...)
}

func Env(key string, defVal interface{}) interface{} {
	return serve.Env(key, defVal)
}

func Object(key string) contracts.Config {
	return serve.Object(key)
}

func String(key string, defVal ...string) string {
	return serve.String(key, defVal...)
}

func Strings(key string, defVal ...string) []string {
	return serve.Strings(key, defVal...)
}

func GoStrings(key string, defVal ...string) goutil.Strings {
	return serve.Strings(key, defVal...)
}

func Int(key string, defVal ...int) int {
	return serve.Int(key, defVal...)
}

func Ints(key string) []int {
	return serve.Ints(key)
}

func Int64(key string, defVal ...int64) int64 {
	return serve.Int64(key, defVal...)
}

func Uint(key string, defVal ...uint) uint {
	return serve.Uint(key, defVal...)
}

func Bool(key string, defVal ...bool) bool {
	return serve.Bool(key, defVal...)
}
