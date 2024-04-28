package fastgo

import (
	"context"
	"net"
)

type Context struct {
	context.Context
	Request  *Request
	Response *Response
}

var (
	FContext = new(Context)
)

// ServerBaseContext 服务端基础context
func ServerBaseContext(listener net.Listener) context.Context {
	return FContext
}

// ServerConnContext 当前链接context
func ServerConnContext(ctx context.Context, c net.Conn) context.Context {
	return ctx
}
