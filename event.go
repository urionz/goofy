package goofy

import (
	"github.com/gookit/event"
	"github.com/urionz/goofy/contracts"
)

func (app *Application) AddListeners(eventListeners contracts.EventListeners) contracts.Application {
	for name, listeners := range eventListeners {
		for _, listener := range listeners {
			app.Manager.AddListener(name, listener)
		}
	}
	return app
}

func (app *Application) MustEmit(name string, payload contracts.EventM) contracts.Event {
	return app.Manager.MustFire(name, payload)
}

func (app *Application) Emit(name string, payload contracts.EventM) (error, contracts.Event) {
	return app.Manager.Fire(name, payload)
}

func (app *Application) Dispatch(name string, payload contracts.EventM) contracts.Application {
	app.Manager.AsyncFire(event.NewBasic(name, payload))
	return app
}

func Listeners(listeners ...contracts.Listener) []contracts.Listener {
	return listeners
}
