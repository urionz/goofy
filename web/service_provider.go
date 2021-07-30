package web

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	irisRecover "github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/urionz/color"
	"github.com/urionz/goofy/command"
	"github.com/urionz/goofy/container"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goofy/log"
	"github.com/urionz/goofy/web/middleware"
)

func init() {
	contracts.AddConfTpl(`
# http相关配置
[http]
# web监听端口
port = 4000
# 是否开启接入日志
access_log = true
`)
}

func NewServiceProvider(app contracts.Application, store *filesystem.Manager) error {
	webEngine := iris.New()
	ctn := webEngine.ConfigureContainer()
	ctn.Use(requestid.New(), irisRecover.New(), logger.New(), middleware.InjectWebContext(store))
	server := &Server{
		Application: webEngine,
	}
	app.AddCommanders(server)
	if err := app.ProvideValue(server, container.Tags{"name": "web"}, container.As(new(contracts.DynamicConf))); err != nil {
		return err
	}
	return app.ProvideValue(ctn, container.Tags{"name": "router"})
}

type Server struct {
	name      string
	port      int
	accessLog bool
	*iris.Application
}

func (e *Server) DynamicConf(_ contracts.Application, conf contracts.Config) error {
	e.accessLog = conf.Bool("http.access_log", false)
	e.SetName(conf.String("app.name", e.name))
	e.SetLevel(conf.String("logger.level", "info"))
	e.Logger().SetOutput(io.MultiWriter(log.GetRotateWriter("logger"), os.Stdout))

	if conf.String("app.env", "production") == "production" {
		e.Logger().SetFormat("json")
	}

	return nil
}

func (e *Server) SetLevel(level string) {
	e.Logger().SetLevel(level)
}

func (e *Server) SetName(name string) {
	e.Application.SetName(name)
}

func (e *Server) Handle(app contracts.Application) *command.Command {
	var conf contracts.Config
	if err := app.Resolve(&conf); err != nil {
		panic(err)
	}
	cmd := &command.Command{
		Name:    "web",
		Aliases: []string{"http"},
		Desc:    "开启http服务",
		Config: func(c *command.Command) {
			c.BoolOpt(&e.accessLog, "log", "", conf.Bool("http.access_log", false), "是否开启接入日志")
			c.StrOpt(&e.name, "name", "n", conf.String("app.name", "web"), "web应用名称")
			c.IntOpt(&e.port, "port", "p", conf.Int("http.port", 3000), "web服务监听端口")
		},
		Func: func(c *command.Command, args []string) error {
			addr := fmt.Sprintf("0.0.0.0:%d", e.port)

			e.UseRouter(iris.NewConditionalHandler(func(ctx iris.Context) bool {
				return e.accessLog
			}, makeAccessLog(conf).Handler))

			if err := e.DynamicConf(app, conf); err != nil {
				log.Error(err)
				return err
			}

			color.Println("<green>[INFO] </>",
				fmt.Sprintf(
					"======================== WebEngine (Port: %d, AppName: %s, EnvName: %s) ========================\n",
					e.port, e.name, conf.String("app.env", "production"),
				),
			)

			return e.Listen(
				addr, iris.WithOptimizations,
				iris.WithConfiguration(iris.Configuration{
					DisableStartupLog:      true,
					RemoteAddrHeadersForce: false,
					RemoteAddrHeaders: []string{
						"X-Real-Ip",
						"X-Forwarded-For",
						"CF-Connecting-IP",
						"True-Client-Ip",
					},
				}),
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
