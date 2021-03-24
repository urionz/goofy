package contracts

import (
	"github.com/goava/di"
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
	AddSchedules(scheduleJob ScheduleJob) Application
	AddServices(services ...interface{}) Application
	AddCommanders(commander ...Commander) Application
	AddListeners(eventListeners EventListeners) Application
	Dispatch(name string, payload EventM) Application
	MustEmit(name string, payload EventM) Event
	Emit(name string, payload EventM) (error, Event)
	ProvideValue(value di.Value, options ...di.ProvideOption) error
	Provide(constructor di.Constructor, options ...di.ProvideOption) error
	Resolve(ptr di.Pointer, options ...di.ResolveOption) error
	Call(args ...string) int

	Workspace() string
	Storage() string
	Database() string
	Run() Application
	Error() error
}
