package main

import (
	"fmt"
	"mime/multipart"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/log"
	"github.com/urionz/goofy/web"
	"github.com/urionz/goofy/web/context"
	"github.com/urionz/goofy/web/validation"
)

func main() {
	goofy.Default.AddServices(
		func(router *iris.APIContainer, app contracts.Application) {
			router.PartyFunc("/", func(router *iris.APIContainer) {
				router.PartyFunc("/idol", func(route *iris.APIContainer) {
					mvc.New(route.Self).Handle(new(Test))
				})
			})
		},
	).Run()
}

type Test struct {
}

type Req struct {
	Files    []*multipart.FileHeader `valid:"mime(image/png)~错误" form:"file[]"`
	Name     string                  `valid:"optional~名称不存在" form:"name"`
	Nickname string                  `form:"nickname"`
}

func (*Test) Post(ctx *context.Context, validate *validation.Validation) *web.JsonResult {
	var req Req
	if err := validate.Validate(ctx, &req); err != nil {
		log.Error(err)
		return web.JsonError(err)
	}
	log.Error("test error")
	log.Infof("hahaha test infgo")
	fmt.Println(req.Files)
	return web.JsonSuccess()
}

func TestHandler(c iris.Context) {
	// fmt.Println(cache.Store().Put("1", "1", 1*time.Minute))
	// fmt.Println(cache.Store().Get("1"))
	// fmt.Println(cache.Store().Forget("1"))
	//
	// fmt.Println(cache.Store().Get("1"))
}
