package db

import (
	"github.com/goava/di"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/db/migrate"
	"github.com/urionz/goofy/db/model"
	"github.com/urionz/goofy/db/seed"
)

func init() {
	contracts.AddConfTpl(`
# 数据库配置
[database]
# 日志级别
log_level = "error"
# 默认数据库连接
default = "test"
	# 数据库连接配置
	[database.conns.test]
	# 用户名
	user = "root"
	# 密码
	password = "root"
	# ip地址
	host = "localhost"
	# 端口
	port = 3306
	# 数据库名称
	name = "test"
	# 字符集
	charset = "utf8mb4"
	# 前缀
	prefix = ""
	# 迁移自动生成外键约束
	auto_migrate_constraint = false
	# 是否开启表名单数
	singular_table = false
	# 慢日志阈值（毫秒）
	slow_threshold = 100
	# 最大连接数
	max_open_conns = 10
	# 最大idle数
	max_idle_conns = 10
`)
}

func NewServiceProvider(app contracts.Application, conf contracts.Config) error {
	app.AddCommanders(
		contracts.FuncCommander(migrate.Make), contracts.FuncCommander(migrate.Migrate),
		contracts.FuncCommander(migrate.Rollback), contracts.FuncCommander(migrate.Status),
		contracts.FuncCommander(migrate.Reset), contracts.FuncCommander(migrate.Refresh),
		contracts.FuncCommander(migrate.Fresh), contracts.FuncCommander(model.Make),
		contracts.FuncCommander(model.MakeRepo), contracts.FuncCommander(seed.Seed),
		contracts.FuncCommander(model.MakeService),
	)
	instance = NewManager(conf)
	return app.ProvideValue(instance, di.As(new(contracts.DBFactory), new(contracts.DynamicConf)))
}
