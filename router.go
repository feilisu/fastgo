package fastgo

import (
	"strings"
)

const POST = "POST"
const GET = "GET"
const PUT = "PUT"
const DELETE = "DELETE"
const HEAD = "HEAD"
const OPTION = "OPTION"
const ANY = "*"

var DefaultHost = "default"

type Handler interface {
	Handle(ctx *Context) error
}

type HandlerFunc func(ctx *Context) error

func (hf HandlerFunc) Handle(ctx *Context) error {
	return hf(ctx)
}

type Router struct {
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

	host := DefaultHost
	if len(f.entry.host) > 0 {
		host = f.entry.host
	}
	if f.entryMap == nil {
		f.entryMap = make(map[string][]Entry)
	}
	f.entryMap[host] = append(f.entryMap[host], f.entry)
	f.entry = Entry{}

}

func (f *Router) GET(handle func(ctx *Context) error) {
	f.method([]string{GET}, HandlerFunc(handle))
}

func (f *Router) POST(handle func(ctx *Context) error) {
	f.method([]string{POST}, HandlerFunc(handle))
}

func (f *Router) DELETE(handle func(ctx *Context) error) {
	f.method([]string{DELETE}, HandlerFunc(handle))
}

func (f *Router) PUT(handle func(ctx *Context) error) {
	f.method([]string{PUT}, HandlerFunc(handle))
}

func (f *Router) HEAD(handle func(ctx *Context) error) {
	f.method([]string{HEAD}, HandlerFunc(handle))
}

func (f *Router) ANY(handle func(ctx *Context) error) {
	f.method([]string{ANY}, HandlerFunc(handle))
}
