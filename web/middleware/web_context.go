package middleware

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy/filesystem"
	"github.com/urionz/goofy/web/context"
	"github.com/urionz/goofy/web/storage"
	"github.com/urionz/goofy/web/validation"
)

func InjectWebContext(manager *filesystem.Manager) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		store := storage.NewStorage(manager)
		fmt.Println("store === ", store)
		depends := []interface{}{
			&context.Context{
				Context: ctx,
				Storage: store,
			},
			&validation.Validation{},
		}

		for _, dep := range depends {
			ctx.RegisterDependency(dep)
		}
	}
}
