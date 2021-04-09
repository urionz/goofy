package contracts

import (
	"github.com/goava/di"
	"github.com/gookit/event"
	"github.com/gookit/gcli/v3"
	"github.com/kataras/iris/v12"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	DITags   = di.Tags
	DIOption = di.Option

	Command = *gcli.Command

	LogField = zap.Field

	DB = *gorm.DB

	Event          = event.Event
	EventM         = event.M
	Listener       = event.Listener
	ListenerFunc   = event.ListenerFunc
	EventListeners = map[string][]Listener

	ScheduleJob = map[string][]interface{}
	FuncJob     = cron.FuncJob

	WebEngine = *iris.Application
	Router    = iris.APIContainer
	Group     = iris.Party
)
