package main

import (
	"fmt"
	"github.com/jinouy/msgo"
	"log"
	"net/http"
)

func Log(next msgo.HandlerFunc) msgo.HandlerFunc {
	return func(ctx *msgo.Context) {
		fmt.Println("打印请求函数")
		next(ctx)
		fmt.Println("返回执行时间")
	}
}

type User struct {
	Name      string   `xml:"name" json:"name"`
	Age       int      `xml:"age" json:"age" validate:"required,max=50,min=18"`
	Addresses []string `json:"addresses"`
	Email     string   `json:"email" msgo:"required"`
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

	//g.Get("/*/get", func(ctx *msgo.Context) {
	//	fmt.Println("handler")
	//	fmt.Fprintln(ctx.W, " get hello joy.com")
	//}, Log)
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

	g.Get("/html", func(ctx *msgo.Context) {
		ctx.HTML(http.StatusOK, "<h1>msgo</h1>")
	})
	g.Get("/htmlTemplate", func(ctx *msgo.Context) {
		ctx.HTMLTemplate("index.html", "", "tpl/index.html")
	})

	engine.LoadTemplate("tpl/*.html")
	g.Get("/template", func(ctx *msgo.Context) {
		user := &User{
			Name: "joy",
		}
		err := ctx.Template("login.html", user)
		if err != nil {
			log.Panicln(err)
		}
	})

	g.Get("/json", func(ctx *msgo.Context) {
		user := &User{
			Name: "joy",
		}
		err := ctx.JSON(http.StatusOK, user)
		if err != nil {
			log.Panicln(err)
		}
	})

	g.Get("/xml", func(ctx *msgo.Context) {
		user := &User{
			Name: "joy",
		}
		err := ctx.XML(http.StatusOK, user)
		if err != nil {
			log.Panicln(err)
		}
	})

	g.Get("/excel", func(ctx *msgo.Context) {
		ctx.File("tpl/test.xlsx")
	})

	g.Get("/excelName", func(ctx *msgo.Context) {
		ctx.FileAttachment("tpl/test.xlsx", "aaaa.xlsx")
	})

	g.Get("/fs", func(ctx *msgo.Context) {
		ctx.FileFromFS("test.xlsx", http.Dir("tpl"))
	})

	g.Get("/redirect", func(ctx *msgo.Context) {
		ctx.Redirect(http.StatusFound, "/user/template")
	})
	g.Get("/string", func(ctx *msgo.Context) {
		ctx.String(http.StatusOK, "%s %s 开始学习如何搭建框架", "从零开始", "joy")
	})

	g.Get("/add", func(ctx *msgo.Context) {
		name := ctx.GetDefaultQuery("name", "张三")
		fmt.Printf("id:%v , ok: %v \n", name, true)
	})
	g.Get("/queryMap", func(ctx *msgo.Context) {
		m, _ := ctx.GetQueryMap("user")
		ctx.JSON(http.StatusOK, m)
	})
	g.Post("/formPost", func(ctx *msgo.Context) {
		m, _ := ctx.GetPostFormMap("user")
		//file := ctx.FormFile("file")
		//err := ctx.SaveUploadFile(file, "./upload/"+file.Filename)
		//if err != nil {
		//	log.Println(err)
		//}
		files := ctx.FormFiles("file")
		for _, file := range files {
			ctx.SaveUploadFile(file, "./upload/"+file.Filename)
		}
		ctx.JSON(http.StatusOK, m)
	})

	g.Post("/jsonParam", func(ctx *msgo.Context) {
		user := make([]User, 0)
		ctx.DisallowUnknownFields = true
		ctx.IsValidate = true
		err := ctx.BindJson(&user)
		if err == nil {
			ctx.JSON(http.StatusOK, user)
		} else {
			log.Println(err)
		}
	})

	g.Post("/xmlParam", func(ctx *msgo.Context) {
		user := &User{}
		err := ctx.BindXML(user)
		log.Println(err)
		ctx.JSON(http.StatusOK, user)
	})
	engine.Run()

}
