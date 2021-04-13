package goofy

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/goava/di"
	"github.com/gookit/event"
	"github.com/gookit/gcli/v3"
	"github.com/robfig/cron/v3"
	"github.com/urionz/goofy/cache"
	"github.com/urionz/goofy/cmds/dlv"
	"github.com/urionz/goofy/cmds/repl"
	"github.com/urionz/goofy/config"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/db"
	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goofy/log"
	"github.com/urionz/goofy/redis"
	"github.com/urionz/goofy/web"
)

var Default = New(SetWorkspace("./")).AddServices(
	config.NewServiceProvider, log.NewServiceProvider,
	redis.NewServiceProvider, db.NewServiceProvider,
	cache.NewServiceProvider, filesystem.NewServiceProvider,
	web.NewServiceProvider,
).AddCommanders(
	contracts.FuncCommander(dlv.Command),
	contracts.FuncCommander(repl.Command),
)

type Application struct {
	di.Tags `name:"app"`

	*cron.Cron
	*event.Manager
	*di.Container
	*gcli.App

	conf      string
	workspace string
	services  []interface{}

	err error
}

var _ contracts.Application = (*Application)(nil)

func New(options ...Option) *Application {
	var err error
	app := &Application{
		workspace: "./",
		Cron:      cron.New(),
		App:       gcli.NewApp(),
		Manager:   event.NewManager("goofy"),
	}

	if app.Container, err = di.New(di.ProvideValue(app, di.As(new(contracts.Application)))); err != nil {
		panic(err)
	}

	app.App.GlobalFlags().StrOpt(&app.workspace, "workspace", "w", "./", "工作目录")

	for _, option := range options {
		option.apply(app)
	}

	return app
}

func (app *Application) Storage() string {
	return path.Join(app.Workspace(), "storage")
}

func (app *Application) Database() string {
	return path.Join(app.Workspace(), "database")
}

func (app *Application) addError(err error) {
	app.err = err
}

func (app *Application) Workspace() string {
	return app.workspace
}

func (app *Application) Dir() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}

func (app *Application) Run() contracts.Application {
	if err := app.Container.Invoke(app.bootstrap); err != nil {
		panic(err)
	}
	if len(os.Args) > 0 && !strings.HasSuffix(os.Args[0], "test") {
		app.App.Run(nil)
	}
	return app
}

func (app *Application) Call(args ...string) int {
	return app.App.Run(args)
}

func (app *Application) Error() error {
	return app.err
}

func (app *Application) bootstrap() error {
	for _, service := range app.services {
		outputs, err := app.resolveInputsFromDI(service)
		if err != nil {
			return err
		}
		if len(outputs) > 0 && !outputs[0].IsNil() {
			return outputs[0].Interface().(error)
		}
	}
	return nil
}

func (app *Application) resolveInputsFromDI(service interface{}) ([]reflect.Value, error) {
	typeOf := reflect.TypeOf(service)
	valueOf := reflect.ValueOf(service)
	serviceArgNum := typeOf.NumIn()
	inputs := make([]reflect.Value, serviceArgNum)
	for i := 0; i < serviceArgNum; i++ {
		newValue := reflect.New(typeOf.In(i))
		ptr := newValue.Interface()
		if err := app.Container.Resolve(ptr); err != nil {
			return []reflect.Value{}, err
		}
		inputs[i] = newValue.Elem()
	}

	if len(os.Args) > 0 && !strings.HasSuffix(os.Args[0], "test") {
		app.App.GlobalFlags().Parse(os.Args[1:])
	}

	return valueOf.Call(inputs), nil
}
