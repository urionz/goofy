package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy/web"
)

func InjectWebContext(ctx iris.Context) {
	depends := []interface{}{
		&web.Context{
			Context: ctx,
		},
		&web.Validation{},

	}

	for _, dep := range depends {
		ctx.RegisterDependency(dep)
	}
}
