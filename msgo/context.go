package msgo

import (
	"errors"
	"github.com/jinouy/msgo/binding"
	"github.com/jinouy/msgo/render"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const defaultMultipartMemory = 32 << 20 //32M
type Context struct {
	W                     http.ResponseWriter
	R                     *http.Request
	engine                *Engine
	queryCache            url.Values
	formCache             url.Values
	DisallowUnknownFields bool
	IsValidate            bool
}

// http://xxx.com/user/add?id=1&age=20&username=张三

func (c *Context) GetDefaultQuery(key, defaultValue string) string {
	values, ok := c.GetQueryArray(key)
	if !ok {
		return defaultValue
	}
	return values[0]
}

func (c *Context) GetQuery(key string) string {
	c.initQueryCache()
	return c.queryCache.Get(key)
}

func (c *Context) QueryArray(key string) (values []string) {
	c.initQueryCache()
	values, _ = c.queryCache[key]
	return values
}

func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.initQueryCache()
	values, ok := c.queryCache[key]
	return values, ok
}

func (c *Context) initQueryCache() {
	if c.R != nil {
		c.queryCache = c.R.URL.Query()
	} else {
		c.queryCache = url.Values{}
	}
}

func (c *Context) QueryMap(key string) (dicts map[string]string) {
	dicts, _ = c.GetQueryMap(key)
	return
}

//http://loaclhost:8080/queryMap?user[id]=1&user[name]=张三
func (c *Context) GetQueryMap(key string) (map[string]string, bool) {
	c.initQueryCache()
	return c.get(c.queryCache, key)
}

func (c *Context) get(cache map[string][]string, key string) (map[string]string, bool) {
	dicts := make(map[string]string)
	exist := false
	//user[id]=1&user[name]=张三
	for k, value := range cache {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
				exist = true
				dicts[k[i+1:][:j]] = value[0]
			}
		}
	}
	return dicts, exist
}

func (c *Context) initPostFormCache() {
	if c.R != nil {
		if err := c.R.ParseMultipartForm(defaultMultipartMemory); err != nil {
			if errors.Is(err, http.ErrNotMultipart) {
				log.Panicln(err)
			}
		}
		c.formCache = c.R.PostForm
	} else {
		c.formCache = url.Values{}
	}
}

func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) PostFormArray(key string) (values []string) {
	values, _ = c.GetPostFormArray(key)
	return
}

func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	c.initPostFormCache()
	values, ok := c.formCache[key]
	return values, ok
}

func (c *Context) PostFormMap(key string) (dicts map[string]string) {
	dicts, _ = c.GetPostFormMap(key)
	return
}

func (c *Context) GetPostFormMap(key string) (map[string]string, bool) {
	c.initPostFormCache()
	return c.get(c.formCache, key)
}

func (c *Context) FormFile(name string) *multipart.FileHeader {
	file, header, err := c.R.FormFile(name)
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()
	return header
}

func (c *Context) FormFiles(name string) []*multipart.FileHeader {
	multipartForm, err := c.MultipartForm()
	if err != nil {
		return make([]*multipart.FileHeader, 0)
	}
	return multipartForm.File[name]
}

func (c *Context) SaveUploadFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.R.ParseMultipartForm(defaultMultipartMemory)
	return c.R.MultipartForm, err
}

func (c *Context) HTML(status int, html string) error {
	//状态是200 默认不设置的话 如果调用了write这个方法，实际上默认返回状态200
	return c.Render(status, &render.HTML{Data: html, IsTemplate: false})
}

func (c *Context) HTMLTemplate(name string, data any, filename ...string) error {
	//状态是200 默认不设置的话 如果调用了write这个方法，实际上默认返回状态200
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.New(name)
	t, err := t.ParseFiles(filename...)
	if err != nil {
		return err
	}
	err = t.Execute(c.W, data)
	return err
}

func (c *Context) HTMLTemplateGlob(name string, data any, pattern string) error {
	//状态是200 默认不设置的话 如果调用了write这个方法，实际上默认返回状态200
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.New(name)
	t, err := t.ParseFiles(pattern)
	if err != nil {
		return err
	}
	err = t.Execute(c.W, data)
	return err
}

func (c *Context) Template(name string, data any) error {
	//状态是200 默认不设置的话 如果调用了write这个方法，实际上默认返回状态200
	return c.Render(http.StatusOK, &render.HTML{
		Data:       data,
		Name:       name,
		Template:   c.engine.HTMLRender.Template,
		IsTemplate: true,
	})
}

func (c *Context) JSON(status int, data any) error {
	//状态是200 默认不设置的话 如果调用了write这个方法，实际上默认返回状态200
	return c.Render(status, &render.JSON{
		Data: data,
	})
}

func (c *Context) XML(status int, data any) error {
	//状态是200 默认不设置的话 如果调用了write这个方法，实际上默认返回状态200
	return c.Render(status, &render.XML{
		Data: data,
	})
}

func (c *Context) File(filename string) {
	http.ServeFile(c.W, c.R, filename)
}

func (c *Context) FileAttachment(filepath, filename string) {
	if isASCII(filename) {
		c.W.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		c.W.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(filepath))
	}
	http.ServeFile(c.W, c.R, filepath)
}

//filepath 是相对文件系统的路径
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.R.URL.Path = old
	}(c.R.URL.Path)

	c.R.URL.Path = filepath

	http.FileServer(fs).ServeHTTP(c.W, c.R)
}

func (c *Context) Redirect(status int, url string) error {
	return c.Render(status, &render.Redirect{
		Code:     status,
		Request:  c.R,
		Location: url,
	})
}

func (c *Context) String(status int, format string, values ...any) error {
	err := c.Render(status, &render.String{Format: format, Data: values})
	c.W.WriteHeader(status)
	return err
}

func (c *Context) Render(statusCode int, r render.Render) error {
	err := r.Render(c.W)
	if statusCode != http.StatusOK {
		c.W.WriteHeader(statusCode)
	}
	return err
}

func (c *Context) BindXML(obj any) error {
	return c.MustBindWith(obj, binding.XML)
}

func (c *Context) BindJson(obj any) error {
	json := binding.JSON
	json.IsValidate = true
	json.DisallowUnknownFields = true
	return c.MustBindWith(obj, json)
}

func (c *Context) MustBindWith(obj any, bind binding.Binding) error {
	if err := c.ShouldBind(obj, bind); err != nil {
		c.W.WriteHeader(http.StatusBadRequest)
		return err
	}
	return nil
}

func (c *Context) ShouldBind(obj any, bind binding.Binding) error {
	return bind.Bind(c.R, obj)

}
