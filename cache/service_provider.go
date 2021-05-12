package cache

import (
	"github.com/urionz/goofy/container"
	"github.com/urionz/goofy/contracts"
)

func init() {
	contracts.AddConfTpl(`
# 缓存配置
[cache]
# 默认缓存配置
default = "file"
# 缓存前缀
prefix = ""
	# 文件缓存
	[cache.stores.file]
	# 驱动名称
	driver = "file"
	# 前缀（优先级大于父级）
	prefix = ""
	# 存储路径
	path = "./"
	
	# redis缓存
	[cache.stores.redis]
	# 驱动名称
	driver = "redis"
	# 前缀（优先级大于父级）
	prefix = ""
	# 连接名称
	connection = "cache"
`)
}

func NewServiceProvider(app contracts.Application, conf contracts.Config) error {
	instance = NewManager(app, conf)
	if err := app.ProvideValue(instance, container.As(new(contracts.CacheFactory))); err != nil {
		return err
	}
	return nil
}
