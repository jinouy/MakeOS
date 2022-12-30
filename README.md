# MakeOS

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