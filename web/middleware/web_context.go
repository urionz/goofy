package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goofy/web/context"
	"github.com/urionz/goofy/web/validation"
)

func InjectWebContext(store *filesystem.Manager) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		depends := []interface{}{
			&context.Context{
				Context: ctx,
				Manager: store,
			},
			&validation.Validation{},
		}

		for _, dep := range depends {
			ctx.RegisterDependency(dep)
		}
	}
}
