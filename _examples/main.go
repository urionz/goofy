package main

import (
	"github.com/kataras/iris/v12"
	"github.com/urionz/goofy"
	_ "github.com/urionz/goofy/_examples/database/migrations"
)

func main() {
	goofy.Default.AddServices(
		func(web *iris.Application) {
			web.Get("/", TestHandler)
		},
	).Run()
}

func TestHandler(c iris.Context) {
	// fmt.Println(cache.Store().Put("1", "1", 1*time.Minute))
	// fmt.Println(cache.Store().Get("1"))
	// fmt.Println(cache.Store().Forget("1"))
	//
	// fmt.Println(cache.Store().Get("1"))
}
