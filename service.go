package goofy

type ServiceFunc func(...interface{}) error

func (app *Application) AddServices(services ...interface{}) IApplication {
	app.services = append(app.services, services...)
	return app
}
