package fastgo

import "net/http"

type Response struct {
	http.ResponseWriter
}

func (r *Response) Text(str string) error {
	r.Header().Set("Content-Type", "text")
	_, err := r.Write([]byte(str))
	if err != nil {
		return err
	}
	return nil
}
