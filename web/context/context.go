package context

import (
	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy/filesystem"
)

type Context struct {
	iris.Context
	*filesystem.Manager
}
