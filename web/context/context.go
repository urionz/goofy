package context

import (
	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy/web/storage"
)

type Context struct {
	iris.Context
	*storage.Storage
}
