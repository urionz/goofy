package web

import "github.com/kataras/iris/v12"

type Context struct {
	iris.Context
}

type Authorizer interface {
	Id() uint
	Username() string
	User() interface{}
}

const (
	AuthContextKey = "auth.context"
	AuthUserId     = "auth.user.id.context"
)

type GetUserFunc func() Authorizer

func NewContext(ctx iris.Context) *Context {
	return &Context{
		Context: ctx,
	}
}

func (ctx *Context) SetUser(f GetUserFunc) *Context {
	ctx.Values().Set(AuthContextKey, f)
	return ctx
}

func (ctx *Context) SetUserId(uid interface{}) *Context {
	ctx.Values().Set(AuthUserId, uid)
	return ctx
}

func (ctx *Context) GetUserId() interface{} {
	return ctx.Values().Get(AuthUserId)
}

func (ctx *Context) User() Authorizer {
	if v := ctx.Values().Get(AuthContextKey); v != nil {
		if f, ok := v.(GetUserFunc); ok {
			return f()
		}
	}
	return nil
}
