package contracts

import (
	"context"
	"github.com/goava/di"
	"github.com/urionz/goofy/event"
	"github.com/urionz/goofy/schedule"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
	FatalLevel = "fatal"
)

type Application interface {
	AddSchedules(scheduleJob schedule.Job) Application
	AddServices(services ...interface{}) Application
	RegisterServices(services ...interface{}) error
	AddCommanders(commander ...Commander) Application
	AddListeners(eventListeners event.Listeners) Application
	Dispatch(name string, payload event.M, fn ...func(error)) Application
	MustEmit(name string, payload event.M) event.Event
	Emit(name string, payload event.M) (error, event.Event)
	ProvideValue(value di.Value, options ...di.ProvideOption) error
	Provide(constructor di.Constructor, options ...di.ProvideOption) error
	Resolve(ptr di.Pointer, options ...di.ResolveOption) error
	Call(name string, args ...string) error
	DynamicConf(conf Config) error

	Workspace() string
	Dir() string
	Storage() string
	Database() string
	Run() Application
	Error() error
}

type DynamicConf interface {
	DynamicConf(app Application, conf Config) error
}

type Closer interface {
	Close(ctx context.Context) error
}
