package log

import (
	"github.com/urionz/goofy/contracts"
)

func init() {
	contracts.AddConfTpl(`
# 日志配置
[logger]
# 日志等级   debug|info|warn|error|panic|fatal
level = "debug"
# 时间格式
encode_time = "2006-01-02 15:04:05"
# 是否开启颜色模式
color = true
`)
}

func NewServiceProvider(app contracts.Application, conf contracts.Config) error {
	return app.ProvideValue(NewLogger(app, conf))
}
