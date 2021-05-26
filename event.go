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

func (app *Application) Dispatch(name string, payload event.M) contracts.Application {
	app.Manager.AsyncFire(event.NewBasic(name, payload))
	return app
}

func Listeners(listeners ...event.Listener) []event.Listener {
	return listeners
}

func ListenerFunc(fn func(event.Event) error) event.ListenerFunc {
	return func(e event.Event) error {
		defer func() {
			if err := recover(); err != nil {
				log.Sugar(0).Error(err)
			}
		}()
		return fn(e)
	}
}
