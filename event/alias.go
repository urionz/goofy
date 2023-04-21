package event

import (
	"github.com/gookit/event"
)

type (
	Event            = event.Event
	M                = event.M
	Listener         = event.Listener
	Listeners        = map[string][]Listener
	ListenerFunc     = event.ListenerFunc
	BasicEvent       = event.BasicEvent
	ArgsListenerFunc func(e event.Event, args ...interface{}) error
)
