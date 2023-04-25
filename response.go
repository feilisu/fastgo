package fastgo

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

func (r *Response) Text(str string) error {
	r.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := r.Write([]byte(str))
	if err != nil {
		return err
	}
	return nil
}

func (r *Response) Json(res any) error {
	r.Header().Set("Content-Type", "application/json;charset=utf8")
	marshal, err := json.Marshal(res)
	if err != nil {
		return err
	}
	_, err = r.Write(marshal)
	if err != nil {
		return err
	}
	return nil
}
