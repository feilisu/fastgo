package fastgo

type Middleware interface {
	Exec(ctx *Context) error
}

type MiddlewareFunc func(*Context) error

func (f MiddlewareFunc) Exec(ctx *Context) error {
	return f(ctx)
}
