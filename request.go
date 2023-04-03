package fastgo

import (
	"io"
	"net/http"
)

type Request struct {
	*http.Request
}

// getString
func (r *Request) getString(key string) string {
	return key
}

func (r *Request) GetJsonBody() string {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	return string(bytes)
}
