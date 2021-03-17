package web

import (
	"github.com/kataras/iris/v12"
)

type Iris = *IrisContext

type IrisContext struct {
	iris.Context
	*Request
}

func IrisWebContext(ctx iris.Context) Iris {
	return &IrisContext{
		Context: ctx,
		Request: NewRequest(ctx.Request()),
	}
}
