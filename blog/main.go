package main

import (
	"fmt"
	"github.com/jinouy/msgo"
)

func main() {

	engine := msgo.New()
	g := engine.Group("user")
	g.Get("/hello/*/get", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, " get hello joy.com")
	})
	g.Post("/hello", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, "post hello joy.com")
	})
	g.Put("/hello", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, "put hello joy.com")
	})
	g.Delete("/hello", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, "delete hello joy.com")
	})
	g.Post("/info", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, "post info joy.com")
	})
	g.Any("/any", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, "any joy.com")
	})
	g.Get("/get/:id", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, "get id joy.com")
	})
	engine.Run()

}
