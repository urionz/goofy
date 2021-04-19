package web

import (
	"io"
	"os"
	"path"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/urionz/goofy/contracts"
)

func makeAccessLog(root string, conf contracts.Config) *accesslog.AccessLog {
	output := path.Join(root, conf.String("logger.output_path", "logs"))
	w, err := rotatelogs.New(path.Join(output, "web-%Y-%m-%d.log"))
	if err != nil {
		panic(err)
	}
	ac := accesslog.New(io.MultiWriter(w, os.Stdout))
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
		ac.SetFormatter(&accesslog.JSON{Indent: "  "})
	}

	return ac
}
