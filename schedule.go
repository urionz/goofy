package goofy

import (
	"fmt"

	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/schedule"
)

func (app *Application) AddSchedules(scheduleJob schedule.Job) contracts.Application {
	for spec, jobs := range scheduleJob {
		for _, job := range jobs {
			switch j := job.(type) {
			case schedule.FuncJob:
				app.AddJob(spec, j)
				break
			case func():
				app.AddFunc(spec, j)
				break
			default:
				app.addError(fmt.Errorf("this %+v type is not support", j))
				break
			}
		}
	}
	return app
}

func Jobs(jobs ...interface{}) []interface{} {
	return jobs
}
