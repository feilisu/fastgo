package fastgo

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var (
	serveMux = http.NewServeMux()
	server   = &http.Server{
		Addr:              ":http",
		Handler:           serveMux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		ErrorLog:          getLogger(),
	}
)

// 获取log记录器
func getLogger() *log.Logger {
	return log.Default()
}

// StartServer 启动服务
func StartServer() {

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}

const POST = "POST"
const GET = "GET"
const PUT = "PUT"
const DELETE = "DELETE"
const HEAD = "HEAD"
const OPTION = "OPTION"
const ANY = "*"

type RouterHandler func(ctx *Context)

type Router struct {
	router *http.ServeMux
	h      string
	p      string
	ms     []string
	handle RouterHandler
}

func NewFastGoRouter() *Router {
	return new(Router)
}

func (f *Router) Host(h string) *Router {
	if len(f.p) > 0 {
		panic("host should init before init path")
	}
	f.h = h
	return f
}

func (f *Router) Tree(call func(f *Router)) {
	if f.p == "" {
		panic("path should init before init Tree")
	}
	call(f)
}

func (f *Router) Path(p string) *Router {
	if p == "" {
		panic("path hot is empty")
	}

	ps := strings.Split(p, "/")
	for _, s := range ps {
		strings.Trim(s, "/")
		f.p = f.p + "/" + s
	}
	return f
}

func (f *Router) Method(ms []string, handle RouterHandler) {
	if ms == nil {
		panic("method not is nil")
	}
	if len(ms) == 0 {
		panic("method not is empty")
	}
	if handle == nil {
		panic("handle not is nil")
	}
	if len(f.p) <= 0 {
		panic("handle should init after init path")
	}

	f.ms = ms
	f.handle = handle
	f.register()
}

func (f *Router) GET(handle RouterHandler) {
	f.Method([]string{GET}, handle)
}

func (f *Router) POST(handle RouterHandler) {
	f.Method([]string{POST}, handle)
}

func (f *Router) DELETE(handle RouterHandler) {
	f.Method([]string{DELETE}, handle)
}

func (f *Router) PUT(handle RouterHandler) {
	f.Method([]string{PUT}, handle)
}

func (f *Router) HEAD(handle RouterHandler) {
	f.Method([]string{HEAD}, handle)
}

func (f *Router) ANY(handle RouterHandler) {
	f.Method([]string{ANY}, handle)
}

func (f *Router) register() {

	if f.router == nil {
		f.router = serveMux
	}

	f.router.Handle(f.h+f.p, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancelFunc := context.WithTimeout(r.Context(), 30*time.Second)

		go func() {
			run(&Context{
				ctx:      ctx,
				Request:  &Request{r},
				Response: &Response{w},
			}, f.ms, f.handle)

		}()

		select {
		case <-ctx.Done():
			cancelFunc()
		}

	}))
}

func run(goContext *Context, methods []string, handler RouterHandler) {

	defer func() {
		a := recover()
		if a != nil {
			panic(a)
		}
	}()

	//服务中间件执行
	MiddlewareRun(MiddlewareParam{})

	//请求方式校验
	if !validRequestMethods(goContext.Request.Method, methods) {
		panic("[" + strings.Join(methods, ",") + ":" + goContext.Request.URL.String() + "] not support request method " + goContext.Request.Method)
	}

	_ = reflect.ValueOf(handler).Call([]reflect.Value{reflect.ValueOf(goContext)})

	goContext.ctx.Done()
}

// validRequestMethods
func validRequestMethods(m string, allowMs []string) bool {
	var normal bool
	for _, am := range allowMs {
		if m == am {
			normal = true
			break
		}
	}
	return normal
}
