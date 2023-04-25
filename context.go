package fastgo

import (
	"context"
)

type Context struct {
	context.Context
	Middlewares []Middleware
	Request     *Request
	Response    *Response
}

func (c *Context) ExecMiddleware() error {
	for _, m := range c.Middlewares {
		err := m.Exec(c)
		if err != nil {
			return err
		}
	}
	return nil
}
