package main

import (
	"fmt"
	"mime/multipart"

	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/cache"
	"github.com/urionz/goofy/redis"
	"github.com/urionz/goofy/web"
	"github.com/urionz/goofy/web/context"
	"github.com/urionz/goofy/web/validation"
)

func main() {
	goofy.Default.AddServices(
		func(rdm *redis.Manager, cm *cache.Manager) {
			var m string
			fmt.Println(cm.Driver().Sear("testfunc", func() interface{} {
				return "value....."
			}, &m))
			var mm string
			cm.Driver().Scan("testfunc", &mm)
			fmt.Println(m, mm)
			// conn, _ := rdm.Connection()
			// fmt.Println(conn.Set("testset", "testvalue", 0))
			// router.PartyFunc("/", func(router *iris.APIContainer) {
			// 	router.PartyFunc("/idol", func(route *iris.APIContainer) {
			// 		mvc.New(route.Self).Handle(new(Test))
			// 	})
			// })
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
	// var req Req
	// if err := validate.Validate(ctx, &req); err != nil {
	// 	log.Error(err)
	// 	return web.JsonError(err)
	// }

	return web.JsonSuccess()
}

func TestHandler(c iris.Context) {
	// fmt.Println(cache.Store().Put("1", "1", 1*time.Minute))
	// fmt.Println(cache.Store().Get("1"))
	// fmt.Println(cache.Store().Forget("1"))
	//
	// fmt.Println(cache.Store().Get("1"))
}
