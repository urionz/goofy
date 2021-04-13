package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goofy/web/context"
	"github.com/urionz/goofy/web/storage"
	"github.com/urionz/goofy/web/validation"
)

func InjectWebContext(manager *filesystem.Manager) func(ctx iris.Context) {
	validate := validation.NewValidation()
	return func(ctx iris.Context) {
		store := storage.NewStorage(manager)
		depends := []interface{}{
			&context.Context{
				Context: ctx,
				Storage: store,
			},
			validate,
		}

		for _, dep := range depends {
			ctx.RegisterDependency(dep)
		}
	}
}
