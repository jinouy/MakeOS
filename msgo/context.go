package msgo

import (
	"github.com/jinouy/msgo/render"
	"html/template"
	"net/http"
	"net/url"
)

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	engine *Engine
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
	c.W.WriteHeader(statusCode)
	return err
}
