package main

import (
	"fmt"
	"github.com/jinouy/msgo"
)

func Log(next msgo.HandlerFunc) msgo.HandlerFunc {
	return func(ctx *msgo.Context) {
		fmt.Println("打印请求函数")
		next(ctx)
		fmt.Println("返回执行时间")
	}
}

func main() {

	engine := msgo.New()
	g := engine.Group("user")

	g.Use(func(next msgo.HandlerFunc) msgo.HandlerFunc {
		return func(ctx *msgo.Context) {
			fmt.Println("pre handle")
			next(ctx)
			fmt.Println("post handle")
		}
	})

	g.Get("/*/get", func(ctx *msgo.Context) {
		fmt.Println("handler")
		fmt.Fprintln(ctx.W, " get hello joy.com")
	}, Log)
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
