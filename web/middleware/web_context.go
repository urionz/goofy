package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy/web"
)

func InjectWebContext(ctx iris.Context) {
	ctx.RegisterDependency([]interface{}{
		web.NewContext(ctx),
		&web.Validation{},
	})
}
