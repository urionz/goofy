package redis

import (
	"github.com/urionz/goofy/contracts"
)

func init() {
	contracts.AddConfTpl(`
# redis连接配置
[redis]
default = "default"
	# 默认redis连接
	[redis.conns.default]
	# 连接地址
	address = "localhost:6379"
	# 连接密码
	password = ""
	# 数据库
	db = 1
	
	# redis缓存连接
	[redis.conns.cache]
	# 连接地址
	address = "localhost:6379"
	# 连接密码
	password = ""
	# 数据库
	db = 2
`)
}

func NewServiceProvider(app contracts.Application, conf contracts.Config) error {
	instance = NewRedisManager(app, conf)
	return app.ProvideValue(instance)
}
