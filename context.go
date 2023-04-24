package fastgo

import (
	"context"
)

type Context struct {
	context.Context
	Request  *Request
	Response *Response
}
