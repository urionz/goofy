package contracts

import "github.com/gookit/event"

type (
	EventListeners map[string][]Listener
	Event          = event.Event
	EventM         = event.M
	Listener       = event.Listener
	ListenerFunc   = event.ListenerFunc
)
