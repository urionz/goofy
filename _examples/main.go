package main

import (
	"fmt"
	"mime/multipart"
	"reflect"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/urionz/goofy"
	_ "github.com/urionz/goofy/_examples/database/migrations"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/validator"
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
	validation.BaseValidator
	Avatar multipart.File `form:"avatar"`
}

func (*Req) Rules(_ *context.Context) validator.MapData {
	return validator.MapData{
		"file:avatar": []string{"ext:zip", "required"},
	}
}

func (*Test) Post(ctx *context.Context, validate *validation.Validation) {
	fmt.Println(reflect.TypeOf(ctx.Manager.Disk()), "test")
	// var r Req
	// if err := validate.Validate(ctx, &r); err != nil {
	// 	log.Error(err)
	// }
	// fmt.Println(r.Avatar)
}

func TestHandler(c iris.Context) {
	// fmt.Println(cache.Store().Put("1", "1", 1*time.Minute))
	// fmt.Println(cache.Store().Get("1"))
	// fmt.Println(cache.Store().Forget("1"))
	//
	// fmt.Println(cache.Store().Get("1"))
}
