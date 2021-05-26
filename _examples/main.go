package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/event"
	"github.com/urionz/goofy/web"
	"github.com/urionz/goofy/web/context"
	"github.com/urionz/goofy/web/validation"
)

func main() {
	goofy.Default.AddListeners(event.Listeners{
		"aaa": []event.Listener{
			goofy.ListenerFunc(func(e event.Event) error {
				return nil
			}),
		},
	}).AddServices(
		func(r *web.Server) {
			mvc.New(r.Application).Handle(new(Test))
		},
	).Run()
}

type Test struct {
}

type Req struct {
	Title    string   `json:"title"`
	Content  string   `valid:"required~必须填写内容" json:"content"`
	Pictures []string `valid:"optional" json:"pictures"`
	Pid      *uint    `valid:"optional" json:"pid"`
	IdolId   uint     `valid:"required_without(Pid)~必须关联idol" json:"idol_id"`
	TopicId  uint     `valid:"optional" json:"topic_id"`
}

func (*Test) Post(ctx *context.Context, validate *validation.Validation) *web.JsonResult {
	var req Req
	if err := validate.Validate(ctx, &req); err != nil {
		return web.JsonError(err)
	}

	return web.JsonSuccess()
}

func TestHandler(c iris.Context) {
	// fmt.Println(cache.Store().Put("1", "1", 1*time.Minute))
	// fmt.Println(cache.Store().Get("1"))
	// fmt.Println(cache.Store().Forget("1"))
	//
	// fmt.Println(cache.Store().Get("1"))
}
