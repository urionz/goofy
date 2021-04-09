package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy/web"
	"github.com/urionz/goofy/web/context"
)

func InjectWebContext(ctx iris.Context) {
	depends := []interface{}{
		&context.Context{
			Context: ctx,
		},
		&web.Validation{},
	}

	for _, dep := range depends {
		ctx.RegisterDependency(dep)
	}
}
