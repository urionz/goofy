package goofy

import (
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/event"
	"github.com/urionz/goofy/log"
)

func (app *Application) AddListeners(eventListeners event.Listeners) contracts.Application {
	for name, listeners := range eventListeners {
		for _, listener := range listeners {
			app.Manager.AddListener(name, listener)
		}
	}
	return app
}

func (app *Application) MustEmit(name string, payload event.M) event.Event {
	return app.Manager.MustFire(name, payload)
}

func (app *Application) Emit(name string, payload event.M) (error, event.Event) {
	return app.Manager.Fire(name, payload)
}

func (app *Application) Dispatch(name string, payload event.M, fn ...func(e error)) contracts.Application {
	go func() {
		err, _ := app.Manager.Fire(name, payload)
		if len(fn) > 0 {
			fn[0](err)
		}
	}()
	return app
}

func Listeners(listeners ...event.Listener) []event.Listener {
	return listeners
}

func ListenerFunc(fn func(event.Event) error) event.ListenerFunc {
	return func(e event.Event) error {
		defer func() {
			if err := recover(); err != nil {
				log.Warn(err)
			}
		}()
		return fn(e)
	}
}
