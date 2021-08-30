package main

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy"
	"github.com/urionz/goofy/contracts"
	"github.com/urionz/goofy/event"
	"github.com/urionz/goofy/web"
	"github.com/urionz/goofy/web/context"
	"github.com/urionz/goofy/web/validation"
)

func main() {
	// fd, _ := os.Open("./config.dev.toml")
	// b, _ := ioutil.ReadAll(fd)
	// fmt.Println(string(b))
	// fmt.Println(string([]byte(`mobiles = ["1", "2", "3", "4"]`)))
	goofy.Default.AddServices(func(conf contracts.Config) {
		fmt.Println(conf.Object("debug").Strings("mobiles"))
	}).Run()
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
	goofy.Default.Emit("aaa", event.M{})
	return web.JsonSuccess()
}

func TestHandler(c iris.Context) {
	// fmt.Println(cache.Store().Put("1", "1", 1*time.Minute))
	// fmt.Println(cache.Store().Get("1"))
	// fmt.Println(cache.Store().Forget("1"))
	//
	// fmt.Println(cache.Store().Get("1"))
}
