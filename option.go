package goofy

type option func(app *Application)

func (o option) apply(app *Application) { o(app) }

type Option interface {
	apply(app *Application)
}

func SetWorkspace(workspace string) Option {
	return option(func(app *Application) {
		app.workspace = workspace
	})
}
