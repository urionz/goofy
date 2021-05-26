package event

import (
	"github.com/gookit/event"
	"github.com/urionz/goofy/log"
)

type (
	Event      = event.Event
	M          = event.M
	Listener   = event.Listener
	Listeners  = map[string][]Listener
	BasicEvent = event.BasicEvent
)

func ListenerFunc(fn func(event.Event) error) event.ListenerFunc {
	return func(e event.Event) error {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()
		return fn(e)
	}
}
