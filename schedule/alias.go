package schedule

import "github.com/robfig/cron/v3"

type (
	Job     = map[string][]interface{}
	FuncJob = cron.FuncJob
)
