package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
	"github.com/urionz/goofy"
	_ "github.com/urionz/goofy/_examples/database/migrations"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/web"
	"github.com/urionz/goofy/web/context"
	"github.com/urionz/goofy/web/middleware"
)

func main() {
	goofy.Default.AddServices(
		func(engine *iris.Application, app contracts.Application) {
			c := engine.ConfigureContainer()
			c.Use(middleware.InjectWebContext)
			c.PartyFunc("/api", func(route *router.APIContainer) {
				route.PartyFunc("/idol", func(route *router.APIContainer) {
					mvc.New(route.Self).Handle(new(Test))
				})
			})
		},
	).Run()
}

type Test struct {
}

type Req struct {
	web.BaseValidator
}

func (*Test) Get(ctx *context.Context, validate *web.Validation) {
	var r Req
	validate.Validate(ctx, &r)
}

func TestHandler(c iris.Context) {
	// fmt.Println(cache.Store().Put("1", "1", 1*time.Minute))
	// fmt.Println(cache.Store().Get("1"))
	// fmt.Println(cache.Store().Forget("1"))
	//
	// fmt.Println(cache.Store().Get("1"))
}
