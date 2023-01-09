
# 参数处理
## 1. query参数
> 首先处理query参数，比如：http://xxx.com/user/add?id=1&age=20&username=张三
> 
> 记得将路由的URL匹配改为：uri := ctx.R.URL.Path

### 1.1 map类型参数
类似于 `http://loaclhost:8080/queryMap?user[id]=1&user[name]=张三`

## 2. Post表单参数
> 获取表单参数借助 http.Request.PostForm

Form属性包含了post表单和url后面跟的get参数。
PostForm属性只包含了post表单参数。

## 3. 文件参数
> 借助http.Request.FormFile

# 错误处理
>当程序发生异常的时候，比如panic程序直接就崩溃了，很明显在应用运行的过程中，不允许发生这样的事情，那么我们的框架就需要支持对这样问题的处理

对这样问题的处理，我们很容易想到recover，他会捕获panic错误，同时我们应当在defer进行捕获，defer是一个先进后出的栈结构，在return之前执行

<b>recover函数只在defer中生效</b>

## 1. Recovery中间件

## 2. 打印出错位置
>像这种处理异常的行为，单纯的报error没有意义，我们需要打印出出错位置

这里我们要用到runtime.Caller, Caller()报告当前go调用栈所执行的函数的文件和行号信息
* `skip`

   上溯的栈帧数，0表示Caller的调用者(Caller所在的调用栈) (0-当前函数,1-上一层函数,...)

* `pc`

   调用栈标识符

* `file`

   文件路径

* `line`

   该调用在文件中的行号

* `ok`