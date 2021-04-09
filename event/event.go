package event

import "github.com/gookit/event"

func NewBasic(name string, data M) *BasicEvent {
	return event.NewBasic(name, data)
}
