package goofy

import "github.com/urionz/goofy/contracts"

func (app *Application) AddServices(services ...interface{}) contracts.Application {
	app.services = append(app.services, services...)
	return app
}

func (app *Application) RegisterServices(services ...interface{}) error {
	for _, service := range services {
		outputs, err := app.resolveInputsFromDI(service)
		if err != nil {
			return err
		}
		if len(outputs) > 0 && !outputs[0].IsNil() {
			return outputs[0].Interface().(error)
		}
	}
	return nil
}
