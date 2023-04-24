package fastgo

import (
	"context"
	"net/http"
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

type Router struct {
	serveMux *http.ServeMux
	entry    Entry
	entryMap map[string][]Entry
}

type Entry struct {
	h      string
	p      string
	ms     []string
	handle RouterHandler
}

func NewRouter() *Router {
	return new(Router)
}

type RouterHandler func(ctx *Context)

func (f *Router) Host(h string) *Router {
	if len(f.entry.p) > 0 {
		panic("host should init before init path")
	}
	f.entry.h = h
	return f
}

func (f *Router) Tree(call func(f *Router)) {
	if f.entry.p == "" {
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
		f.entry.p = f.entry.p + "/" + s
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
	if len(f.entry.p) <= 0 {
		panic("handle should init after init path")
	}

	f.entry.ms = ms
	f.entry.handle = handle

	host := defaultHost
	if len(f.entry.h) > 0 {
		host = f.entry.h
	}
	if f.entryMap == nil {
		f.entryMap = make(map[string][]Entry)
	}
	f.entryMap[host] = append(f.entryMap[host], f.entry)
	f.entry = Entry{}

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

	f.serveMux = http.NewServeMux()
	if f.entryMap != nil {
		for host, entryList := range f.entryMap {
			serveMux := http.NewServeMux()
			for _, entry := range entryList {
				serveMux.Handle(entry.p, getHandler(entry))
			}
			f.serveMux.Handle(host, serveMux)
		}
	}
}

func getHandler(entry Entry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelFunc()
		fastGoContext := &Context{
			Context:  ctx,
			Request:  &Request{r},
			Response: &Response{w},
		}
		entry.handle(fastGoContext)
	})
}
