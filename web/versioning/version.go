package versioning

import (
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/versioning"
	"github.com/urionz/goofy/web/context"
)

type (
	API   = router.Party
	Group = *versioning.Group
)

func NewGroup(r API, version string) Group {
	return versioning.NewGroup(r, version)
}

func SetVersion(ctx *context.Context, constraint string) {
	versioning.SetVersion(ctx.Context, constraint)
}
