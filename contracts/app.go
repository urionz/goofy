package contracts

import (
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
	AddCommanders(commander ...Commander) Application
	AddListeners(eventListeners event.Listeners) Application
	Dispatch(name string, payload event.M) Application
	MustEmit(name string, payload event.M) event.Event
	Emit(name string, payload event.M) (error, event.Event)
	ProvideValue(value di.Value, options ...di.ProvideOption) error
	Provide(constructor di.Constructor, options ...di.ProvideOption) error
	Resolve(ptr di.Pointer, options ...di.ResolveOption) error
	Call(name string, args ...string) error

	Workspace() string
	Dir() string
	Storage() string
	Database() string
	Run() Application
	Error() error
}
