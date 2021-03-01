package goofy

import (
	"bytes"
	"reflect"

	"github.com/goava/di"
	"github.com/gookit/event"
	"github.com/robfig/cron/v3"
	"github.com/urionz/cobra"
)

var Default = New(SetWorkspace("./"))

type (
	IApplication interface {
		AddSchedules(scheduleJob ScheduleJob) IApplication
		AddServices(services ...interface{}) IApplication
		AddCommanders(commander ...Commander) IApplication
		AddListeners(eventListeners EventListeners) IApplication
		Dispatch(name string, payload EventM) IApplication
		MustEmit(name string, payload EventM) Event
		Emit(name string, payload EventM) (error, Event)
		Provide(constructor di.Constructor, options ...di.ProvideOption) error
		Resolve(ptr di.Pointer, options ...di.ResolveOption) error
		Call(args ...string) (*cobra.Command, string, error)

		Workspace() string
		Run() IApplication
		Error() error
	}

	Application struct {
		*cron.Cron
		*event.Manager
		*di.Container
		*cobra.Command

		workspace string
		services  []interface{}

		err error
	}
)

func New(options ...Option) IApplication {
	var err error
	app := &Application{
		workspace: "./",
		Cron:      cron.New(),
		Command: &cobra.Command{
			Use: "goofy",
		},
		Manager: event.NewManager("goofy"),
	}

	app.Command.AddCommand()

	app.Command.InheritedFlags().StringVarP(&app.workspace, "workspace", "w", "./", "执行工作目录")

	if app.Container, err = di.New(di.Provide(func() *Application { return app }, di.As(new(IApplication)))); err != nil {
		panic(err)
	}

	for _, option := range options {
		option.apply(app)
	}

	return app
}

func (app *Application) addError(err error) {
	app.err = err
}

func (app *Application) Workspace() string {
	return app.workspace
}

func (app *Application) Run() IApplication {
	err := app.Container.Invoke(app.run)
	if err != nil {
		panic(err)
	}
	return app
}

func (app *Application) Call(args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	app.Command.SetOut(buf)
	app.Command.SetErr(buf)
	app.Command.SetArgs(args)

	c, err = app.Command.ExecuteC()

	return c, buf.String(), err
}

func (app *Application) Error() error {
	return app.err
}

func (app *Application) run() error {
	if err := app.bootstrap(); err != nil {
		return err
	}
	app.Command.Execute()
	return nil
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
		if typeOf.In(i).Implements(reflect.TypeOf((*IApplication)(nil)).Elem()) {
			inputs[i] = reflect.ValueOf(app)
		} else {
			newValue := reflect.New(typeOf.In(i))
			ptr := newValue.Interface()
			if err := app.Container.Resolve(ptr); err != nil {
				return []reflect.Value{}, err
			}
			inputs[i] = newValue.Elem()
		}
	}
	return valueOf.Call(inputs), nil
}
