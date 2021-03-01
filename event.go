package goofy

import "github.com/gookit/event"

type (
	EventListeners map[string][]Listener
	Event          = event.Event
	EventM         = event.M
	Listener       = event.Listener
	ListenerFunc   = event.ListenerFunc
)

func (app *Application) AddListeners(eventListeners EventListeners) IApplication {
	for name, listeners := range eventListeners {
		for _, listener := range listeners {
			app.Manager.AddListener(name, listener)
		}
	}
	return app
}

func (app *Application) MustEmit(name string, payload EventM) Event {
	return app.Manager.MustFire(name, payload)
}

func (app *Application) Emit(name string, payload EventM) (error, Event) {
	return app.Manager.Fire(name, payload)
}

func (app *Application) Dispatch(name string, payload EventM) IApplication {
	app.Manager.AsyncFire(event.NewBasic(name, payload))
	return app
}

func Listeners(listeners ...Listener) []Listener {
	return listeners
}
