package goofy

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

type (
	ScheduleJob map[string][]interface{}
	FuncJob     = cron.FuncJob
)

func (app *Application) AddSchedules(scheduleJob ScheduleJob) IApplication {
	for spec, jobs := range scheduleJob {
		for _, job := range jobs {
			switch j := job.(type) {
			case cron.FuncJob:
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
