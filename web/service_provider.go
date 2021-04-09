package web

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	irisRecover "github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/urionz/goofy/command"
	"github.com/urionz/goofy/container"
	"github.com/urionz/goofy/contracts"
)

func init() {
	contracts.AddConfTpl(`
# http相关配置
[http]
# web监听端口
port = 4000
`)
}

func NewServiceProvider(app contracts.Application) error {
	webEngine := iris.New()
	webEngine.Use(requestid.New(), irisRecover.New(), logger.New())
	app.AddCommanders(&engine{
		Engine: webEngine,
	})
	return app.ProvideValue(webEngine, container.Tags{"name": "web"})
}

type engine struct {
	name  string
	debug bool
	port  int
	*Engine
}

func (e *engine) Handle(app contracts.Application) *command.Command {
	var conf contracts.Config
	if err := app.Resolve(&conf); err != nil {
		panic(err)
	}
	cmd := &command.Command{
		Name:    "web",
		Aliases: []string{"http"},
		Desc:    "开启http服务",
		Config: func(c *command.Command) {
			c.StrOpt(&e.name, "name", "n", conf.String("app.name", "web"), "web应用名称")
			c.BoolOpt(&e.debug, "debug", "d", conf.Bool("app.debug", false), "是否开启web调试")
			c.IntOpt(&e.port, "port", "p", conf.Int("http.port", 3000), "web服务监听端口")
		},
		Func: func(c *command.Command, args []string) error {
			addr := fmt.Sprintf("0.0.0.0:%d", e.port)
			e.SetName(e.name)
			if e.debug {
				e.Logger().SetLevel("debug")
				e.Logger().SetOutput(os.Stdout)
			}
			return e.Listen(
				addr, iris.WithOptimizations,
				iris.WithRemoteAddrPrivateSubnet("192.168.0.0", "192.168.255.255"),
				iris.WithRemoteAddrPrivateSubnet("10.0.0.0", "10.255.255.255"),
			)
		},
	}
	return cmd
}

func readAppDirectories(directory string, paths *[]string) {
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}

	useDirectory := false
	for _, fileInfo := range fileInfos {
		if useDirectory {
			continue
		}

		if fileInfo.IsDir() && fileInfo.Name()[0] != '.' {
			readAppDirectories(directory+"/"+fileInfo.Name(), paths)
			continue
		}

		if path.Ext(fileInfo.Name()) == ".go" {
			*paths = append(*paths, directory)
			useDirectory = true
		}
	}
}
