package context

import "github.com/kataras/iris/v12/context"

type (
	IrisHandler = context.Handler
	Handler     = func(ctx *Context)
)
