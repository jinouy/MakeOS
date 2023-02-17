package main

import (
	"errors"
	"fmt"
	"github.com/jinouy/msgo"
	msLog "github.com/jinouy/msgo/log"
	"github.com/jinouy/msgo/mserror"
	"github.com/jinouy/msgo/mspool"
	"github.com/jinouy/msgo/token"
	"log"
	"net/http"
	"sync"
	"time"
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

	engine := msgo.Default()
	engine.RegisterErrorHandler(func(err error) (int, any) {
		switch e := err.(type) {
		case *BlogResponse:
			return http.StatusOK, e.Response()
		default:
			return http.StatusInternalServerError, "500 error"
		}
	})
	//fmt.Println(msgo.BasicAuth("joy", "123456"))
	//auth := &msgo.Accounts{
	//	Users: make(map[string]string),
	//}
	//auth.Users["joy"] = "123456"
	//engine.Use(auth.BasicAuth)
	jh := &token.JwtHandler{Key: []byte("123456")}
	//为特定的中间件，需要指定不进行拦截的请求
	engine.Use(jh.AuthInterceptor)
	g := engine.Group("user")
	//g.Use(msgo.Logging, msgo.Recovery)

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

	engine.Logger.Level = msLog.LevelDebug
	//engine.Logger.Formatter = &msLog.JsonFormatter{TimeDisplay: true}
	//logger.Outs = append(logger.Outs, msLog.FileWriter("./log/log.log"))
	engine.Logger.LogFileSize = 1 << 10

	g.Post("/xmlParam", func(ctx *msgo.Context) {
		user := &User{}
		_ = ctx.BindXML(user)
		//ctx.Logger.WithFields(msLog.Fields{
		//	"name": "msgo",
		//	"id":   1,
		//}).Debug("我是debug日志")
		//
		//ctx.Logger.Info("我是info日志")
		//ctx.Logger.Error("我是error日志")
		//err := mserror.Default()
		//err.Result(func(msError *mserror.MsError) {
		//	ctx.Logger.Info(msError.Error())
		//	ctx.JSON(http.StatusInternalServerError, user)
		//})
		//a(1, err)
		//b(1, err)
		//c(1, err)
		//ctx.JSON(http.StatusOK, user)
		err := login()
		ctx.HandleWithError(http.StatusOK, user, err)
	})
	p, _ := mspool.NewPool(5)
	g.Post("/pool", func(ctx *msgo.Context) {
		currentTime := time.Now().UnixMilli()
		var wg sync.WaitGroup
		wg.Add(5)
		p.Submit(func() {
			defer func() {
				wg.Done()
			}()

			fmt.Println("11111111")
			panic("这是1111的panic")
			time.Sleep(3 * time.Second)
		})
		p.Submit(func() {
			fmt.Println("22222222")
			time.Sleep(3 * time.Second)
			wg.Done()
		})
		p.Submit(func() {
			fmt.Println("33333333")
			time.Sleep(3 * time.Second)
			wg.Done()
		})
		p.Submit(func() {
			fmt.Println("44444444")
			time.Sleep(3 * time.Second)
			wg.Done()
		})
		p.Submit(func() {
			fmt.Println("55555555")
			time.Sleep(3 * time.Second)
			wg.Done()
		})
		wg.Wait()
		fmt.Printf("time: %v \n", time.Now().UnixMilli()-currentTime)
		ctx.JSON(http.StatusOK, "success")
	})
	g.Get("/login", func(ctx *msgo.Context) {
		jwt := &token.JwtHandler{}
		jwt.Key = []byte("123456")
		jwt.SendCookie = true
		jwt.TimeOut = 10 * time.Minute
		jwt.RefreshTimeOut = 20 * time.Minute
		jwt.Authenticator = func(ctx *msgo.Context) (map[string]any, error) {
			data := make(map[string]any)
			data["userId"] = 1
			return data, nil
		}
		token, err := jwt.LoginHandler(ctx)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, token)

	})
	g.Get("/refresh", func(ctx *msgo.Context) {
		jwt := &token.JwtHandler{}
		jwt.Key = []byte("123456")
		jwt.SendCookie = true
		jwt.TimeOut = 10 * time.Minute
		jwt.RefreshTimeOut = 20 * time.Minute
		jwt.RefreshKey = []byte("blog_refresh_token")
		ctx.Set(string(jwt.RefreshKey), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzM1MDgxNDIsImlhdCI6MTY3MzUwNjk0MiwidXNlcklkIjoxfQ.H7-0v0CTCDT7j6uuT2PpXw2Gw6qSPyZt4aI7TbdKrXs")
		token, err := jwt.RefreshHandler(ctx)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, token)

	})
	//engine.Run()
	engine.RunTLS(":8118", "key/server.pem", "key/server.key")
}

type BlogResponse struct {
	Success bool
	Code    int
	Data    any
	Msg     string
}

type BlogNoDataResponse struct {
	Success bool
	Code    int
	Msg     string
}

func (b *BlogResponse) Error() string {
	return b.Msg
}

func (b *BlogResponse) Response() any {
	if b.Data == nil {
		return &BlogNoDataResponse{
			Success: false,
			Code:    -999,
			Msg:     "账号密码错误",
		}
	}
	return b
}

func login() *BlogResponse {
	return &BlogResponse{
		Success: false,
		Code:    -999,
		Data:    nil,
		Msg:     "账号密码错误",
	}
}

func a(param int, msError *mserror.MsError) {
	if param == 1 {
		// 发生错误的时候，放入一个地方 然后进行统一处理
		err := errors.New("a error")
		msError.Put(err)
	}
}

func b(param int, msError *mserror.MsError) {
	if param == 1 {

		err := errors.New("a error")
		msError.Put(err)
	}
}

func c(param int, msError *mserror.MsError) {
	if param == 1 {
		err := errors.New("a error")
		msError.Put(err)
	}
}
