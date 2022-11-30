package main

import (
	"fmt"
	"github.com/jinouy/msgo"
)

func main() {

	engine := msgo.New()
	g := engine.Group("user")
	g.Get("/hello", func(ctx *msgo.Context) {
		fmt.Fprintln(ctx.W, "hello joy.com")
	})
	engine.Run()

}
