package gin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

// IResponse代表返回方法
type IResponse interface {
	// IJson 输出
	IJson(obj interface{}) IResponse

	// IJsonp 输出
	IJsonp(obj interface{}) IResponse

	// IXml 输出
	IXml(obj interface{}) IResponse

	// html 输出
	IHtml(template string, obj interface{}) IResponse

	// string
	IText(format string, values ...interface{}) IResponse

	// 重定向
	IRedirect(path string) IResponse

	// header
	ISetHeader(key string, val string) IResponse

	// Cookie
	ISetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	// 设置状态码
	ISetStatus(code int) IResponse

	// 设置 200 状态
	ISetOkStatus() IResponse
}

func (ctx *ContainerContext) IJson(obj interface{}) IResponse {
	byt, err := json.Marshal(obj)
	if err != nil {
		return ctx.ISetStatus(http.StatusInternalServerError)
	}
	ctx.ISetHeader("Content-Type", "application/json")
	ctx.Writer.Write(byt)
	return ctx
}

func (ctx *ContainerContext) IJsonp(obj interface{}) IResponse {
	// 获取请求参数callback
	callbackFunc := ctx.Query("callback")
	ctx.ISetHeader("Content-Type", "application/javascript")
	// 输出到前端页面的时候需要注意下进行字符过滤，否则有可能造成xss攻击
	callback := template.JSEscapeString(callbackFunc)

	// 输出函数名
	_, err := ctx.Writer.Write([]byte(callback))
	if err != nil {
		return ctx
	}
	// 输出左括号
	_, err = ctx.Writer.Write([]byte("("))
	if err != nil {
		return ctx
	}
	// 数据函数参数
	ret, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}
	_, err = ctx.Writer.Write(ret)
	if err != nil {
		return ctx
	}
	// 输出右括号
	_, err = ctx.Writer.Write([]byte(")"))
	if err != nil {
		return ctx
	}
	return ctx
}

func (ctx *ContainerContext) IXml(obj interface{}) IResponse {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return ctx.ISetStatus(http.StatusInternalServerError)
	}
	ctx.ISetHeader("Content-Type", "application/html")
	ctx.Writer.Write(byt)
	return ctx
}

func (ctx *ContainerContext) IHtml(file string, obj interface{}) IResponse {
	// 读取模版文件，创建template实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx
	}
	// 执行Execute方法将obj和模版进行结合
	if err := t.Execute(ctx.Writer, obj); err != nil {
		return ctx
	}

	ctx.ISetHeader("Content-Type", "application/html")
	return ctx
}

func (ctx *ContainerContext) IText(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	ctx.ISetHeader("Content-Type", "application/text")
	ctx.Writer.Write([]byte(out))
	return ctx
}

// 重定向
func (ctx *ContainerContext) IRedirect(path string) IResponse {
	http.Redirect(ctx.Writer, ctx.Request, path, http.StatusMovedPermanently)
	return ctx
}

func (ctx *ContainerContext) ISetHeader(key string, val string) IResponse {
	ctx.Writer.Header().Add(key, val)
	return ctx
}

// Cookie
func (ctx *ContainerContext) ISetCookie(key string, val string, maxAge int, path string, domain string, secure bool, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return ctx
}

// 设置状态码
func (ctx *ContainerContext) ISetStatus(code int) IResponse {
	ctx.Writer.WriteHeader(code)
	return ctx
}

// 设置200状态
func (ctx *ContainerContext) ISetOkStatus() IResponse {
	ctx.Writer.WriteHeader(http.StatusOK)
	return ctx
}
