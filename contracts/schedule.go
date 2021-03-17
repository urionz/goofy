package contracts

import "github.com/robfig/cron/v3"

type (
	ScheduleJob map[string][]interface{}
	FuncJob     = cron.FuncJob
)
