package log

import (
	"github.com/urionz/goofy/contracts"
)

func init() {
	contracts.AddConfTpl(`
# 日志配置
[logger]
# 是否开启日志输出
output_enable = true
# 日志等级   debug|info|warn|error|panic|fatal
level = "debug"
# 时间格式
encode_time = "2006-01-02 15:04:05"
# 是否开启颜色模式
color = true
# 多等级日志输出
multi_level_output = true
# 输出目录
output_path = "logs"
output_path_abs = ""
`)
}

func NewServiceProvider(app contracts.Application, conf contracts.Config) error {
	return app.ProvideValue(NewLogger(app, conf))
}
