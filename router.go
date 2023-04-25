package fastgo

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"
)

const POST = "POST"
const GET = "GET"
const PUT = "PUT"
const DELETE = "DELETE"
const HEAD = "HEAD"
const OPTION = "OPTION"
const ANY = "*"

var defaultHost = "default"

type Handler interface {
	Handle(ctx *Context) error
}

type HandlerFunc func(ctx *Context) error

func (hf HandlerFunc) Handle(ctx *Context) error {
	return hf(ctx)
}

type Router struct {
	serveMux *http.ServeMux
	entry    Entry
	entryMap map[string][]Entry
}

type Entry struct {
	host        string
	path        string
	methods     []string
	middlewares []Middleware
	handle      Handler
}

func NewRouter() *Router {
	return new(Router)
}

func (f *Router) Host(h string) *Router {
	if len(f.entry.path) > 0 {
		panic("host should init before init path")
	}
	f.entry.host = h
	return f
}

func (f *Router) Tree(call func(f *Router)) {
	if f.entry.path == "" {
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
		if len(s) <= 0 {
			continue
		}
		strings.Trim(s, "/")
		f.entry.path = f.entry.path + "/" + s
	}
	return f
}

func (f *Router) Middleware(ms []Middleware) *Router {
	f.entry.middlewares = ms
	return f
}

func (f *Router) method(ms []string, handle Handler) {
	if ms == nil {
		panic("method not is nil")
	}
	if len(ms) == 0 {
		panic("method not is empty")
	}
	if handle == nil {
		panic("handle not is nil")
	}
	if len(f.entry.path) <= 0 {
		panic("handle should init after init path")
	}

	f.entry.methods = ms
	f.entry.handle = handle

	host := defaultHost
	if len(f.entry.host) > 0 {
		host = f.entry.host
	}
	if f.entryMap == nil {
		f.entryMap = make(map[string][]Entry)
	}
	f.entryMap[host] = append(f.entryMap[host], f.entry)
	f.entry = Entry{}

}

func (f *Router) GET(handle Handler) {
	f.method([]string{GET}, handle)
}

func (f *Router) POST(handle Handler) {
	f.method([]string{POST}, handle)
}

func (f *Router) DELETE(handle Handler) {
	f.method([]string{DELETE}, handle)
}

func (f *Router) PUT(handle Handler) {
	f.method([]string{PUT}, handle)
}

func (f *Router) HEAD(handle Handler) {
	f.method([]string{HEAD}, handle)
}

func (f *Router) ANY(handle Handler) {
	f.method([]string{ANY}, handle)
}

func (f *Router) register(ms []Middleware) *http.ServeMux {

	f.serveMux = http.NewServeMux()
	if f.entryMap != nil {
		for host, entryList := range f.entryMap {
			for _, entry := range entryList {
				log.Printf("%s %s %s", entry.methods, host, entry.path)

				var url string
				if entry.host == defaultHost {
					url = entry.path
				} else {
					url = host + entry.path
				}

				f.serveMux.Handle(url, getHandler(entry, ms))
			}
		}
	}

	f.serveMux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	f.serveMux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	f.serveMux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	f.serveMux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))

	return f.serveMux
}

func getHandler(entry Entry, middlewares []Middleware) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelFunc()

		if entry.middlewares != nil {
			middlewares = append(middlewares, entry.middlewares...)
		}

		fctx := &Context{
			Context:     ctx,
			Request:     &Request{r},
			Response:    &Response{w},
			Middlewares: middlewares,
		}

		err := fctx.ExecMiddleware()
		if err != nil {
			_ = fctx.Response.Json(err)
			return
		}

		err = entry.handle.Handle(fctx)
		if err != nil {
			_ = fctx.Response.Json(err)
			return
		}
	})
}
