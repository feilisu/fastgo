package fastgo

import (
	"context"
)

type Context struct {
	ctx      context.Context
	Request  *Request
	Response *Response
}
