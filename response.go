package fastgo

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

func buildResponse(w http.ResponseWriter) *Response {
	return &Response{ResponseWriter: w}
}

const (
	HttpStatus200 = 200
	HttpStatus500 = 500
)

func (r *Response) Text(str string) error {
	bs := []byte(str)
	r.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.WriteHeader(HttpStatus200)
	_, err := r.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (r *Response) Json(res any) error {
	marshal, err := json.Marshal(res)
	if err != nil {
		return err
	}
	r.Header().Set("Access-Control-Allow-Origin", "*")
	r.Header().Set("Content-Type", "application/json;charset=utf8")
	r.WriteHeader(HttpStatus200)
	_, err = r.Write(marshal)
	if err != nil {
		return err
	}
	return nil
}

func (r *Response) Error(str error) error {
	bs := []byte(str.Error())
	r.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.WriteHeader(HttpStatus500)
	_, err := r.Write(bs)
	if err != nil {
		return err
	}
	return nil
}
