package web

import (
	"io"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/log"
)

func makeAccessLog(conf contracts.Config) *accesslog.AccessLog {
	ac := accesslog.New(io.MultiWriter(log.GetRotateWriter("logger"), os.Stdout))
	ac.BytesSent = false
	ac.AddFields(func(ctx iris.Context, fields *accesslog.Fields) {
		fields.Set("request_id", ctx.GetID())
		if ctx.User() != nil {
			if id, err := ctx.User().GetID(); err == nil {
				fields.Set("uid", id)
			}
		}
	})

	if conf.String("app.env", "production") == "production" {
		ac.SetFormatter(&accesslog.JSON{
			HumanTime: true,
		})
	}

	return ac
}
