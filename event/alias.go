package event

import "github.com/gookit/event"

type (
	Event        = event.Event
	M            = event.M
	Listener     = event.Listener
	ListenerFunc = event.ListenerFunc
	Listeners    = map[string][]Listener
	BasicEvent   = event.BasicEvent
)
