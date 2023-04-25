package fastgo

import "log"

type Middleware interface {
	Exec(ctx *Context) error
}

type MiddlewareFunc func(*Context) error

func (f MiddlewareFunc) Exec(ctx *Context) error {
	return f(ctx)
}

func Mtest1(ctx *Context) error {
	log.Print("Mtest1")
	return nil
}

type Mtest2 struct {
}

func (m *Mtest2) Exec(ctx *Context) error {
	log.Print("Mtest2")
	return nil
}
