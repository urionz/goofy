package cmds

import (
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/gcli/v3"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/log"
	"github.com/urionz/goutil/fsutil"
)

var (
	mode = "dev"
	tpl  = `# 应用相关配置
[app]
# 应用名称
name = "app_name"
# 是否开启调试模式
debug = false

# 日志配置
[logger]
# 日志等级   debug|info|warn|error|panic|fatal
level = "debug"
# 时间格式
encode_time = "2006-01-02 15:04:05"
# 是否开启颜色模式
color = true


# http相关配置
[http]
# web监听端口
port = 4000

# 数据库配置
[database]
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
	# 是否开启表名单数
	singular_table = false
	# 慢日志阈值（毫秒）
	slow_threshold = 100
	# 最大连接数
	max_open_conns = 10
	# 最大idle数
	max_idle_conns = 10

	
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
`
)

func GenerateConf(app contracts.Application) *gcli.Command {
	cmd := &gcli.Command{
		Name: "make-conf",
		Desc: "生产环境配置文件",
		Config: func(c *gcli.Command) {
			c.StrOpt(&mode, "mode", "", "dev", "生成环境")
		},
		Func: func(cmd *gcli.Command, args []string) error {
			filePath := path.Join(app.Workspace(), fmt.Sprintf("config.%s.toml", mode))
			if fsutil.FileExists(filePath) {
				cover := false
				prompt := &survey.Confirm{
					Message: "该迁移文件已存在，是否覆盖？",
				}
				survey.AskOne(prompt, &cover)
				if !cover {
					return nil
				}
			}

			if f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666); err == nil {
				f.WriteString(tpl)
				f.Close()
			} else {
				return err
			}

			envPath := path.Join(app.Workspace(), ".env")

			sw := false
			prompt := &survey.Confirm{
				Message: "是否切换当前运行环境？",
			}
			survey.AskOne(prompt, &sw)
			if !sw {
				return nil
			}

			if f, err := os.OpenFile(envPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666); err == nil {
				fmt.Println(f.WriteString(fmt.Sprintf("APP_ENV=%s", mode)))
				f.Close()
			} else {
				log.Error(err)
				return err
			}

			return nil
		},
	}
	return cmd
}
