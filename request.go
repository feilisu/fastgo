package fastgo

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
)

type Request struct {
	*http.Request

	requestParams *requestParams
}

func buildRequest(r *http.Request) (*Request, error) {
	request := &Request{Request: r, requestParams: &requestParams{}}
	err := parseParams(request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

type requestParams struct {
	body io.ReadCloser

	// postParams + getParams; postParams优先级高于getParams, 如：postParams 和 getParams 同有id参数，id的值取postParams
	params map[string]string

	//URL 参数
	getParams map[string]string

	//PATH 参数
	pathParams map[string]string

	//multipart/form-data  application/x-www-form-urlencoded 参数
	postParams map[string]string

	state int
	mux   sync.RWMutex
	files map[string][]*multipart.FileHeader
}

const (
	HEADER_APPLICATION_JSON      = "application/json"
	HEADER_TEXT_HTML             = "text/html"
	HEADER_APPLICATION_XML       = "application/xml"
	HEADER_TEXT_XML              = "text/xml"
	HEADER_TEXT_PLAIN            = "text/plain"
	HEADER_X_WWW_FORM_URLENCODED = "application/x-www-form-urlencoded"
	HEADER_MULTIPART_FORM_DATA   = "multipart/form-data"
)

const defaultMaxMemory int64 = 30 << 20

// getString
func (r *Request) getString(key string) string {
	if res, ok := r.requestParams.params[key]; ok {
		return res
	}

	return ""

}

// GetParams
func (r *Request) GetParams() map[string]string {
	return r.requestParams.params
}

func (r *Request) GetJsonBody() string {
	//r.URL.RawQuery
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (r *Request) QueryParams(param any) (err error) {
	return mapstructure.Decode(r.requestParams.getParams, param)
}

func (r *Request) PathParams(param any) (err error) {
	return mapstructure.Decode(r.requestParams.pathParams, param)

}

func (r *Request) PostParams(param any) (err error) {
	return mapstructure.Decode(r.requestParams.postParams, param)

}

func (r *Request) XmlParams(param any) (err error) {
	decoder := xml.NewDecoder(r.Body)
	err = decoder.Decode(param)
	return
}

func (r *Request) JsonParams(param any) (err error) {
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(param)
	return
}

func (r *Request) Params(param any) (err error) {
	err = mapstructure.Decode(r.requestParams.params, param)
	if err != nil {
		return err
	}

	hv := r.Header.Get("Content-Type")
	if strings.Contains(hv, HEADER_APPLICATION_JSON) {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(param)
	} else if strings.Contains(hv, HEADER_APPLICATION_XML) {
		decoder := xml.NewDecoder(r.Body)
		err = decoder.Decode(param)
	}
	return
}

// parseParams
func parseParams(r *Request) error {
	r.requestParams.mux.Lock()

	if r.requestParams.state == 1 {
		return nil
	}

	hv := r.Header.Get("Content-Type")
	if strings.Contains(hv, HEADER_MULTIPART_FORM_DATA) {
		if err := r.ParseMultipartForm(defaultMaxMemory); err != nil {
			return err
		}
		r.requestParams.postParams = vsTov(r.MultipartForm.Value)
		if r.MultipartForm != nil && r.MultipartForm.File != nil {
			r.requestParams.files = r.MultipartForm.File
		}
	} else if strings.Contains(hv, HEADER_X_WWW_FORM_URLENCODED) {
		if err := r.ParseForm(); err != nil {
			return err
		}
		r.requestParams.postParams = vsTov(r.PostForm)
	} else if strings.Contains(hv, HEADER_APPLICATION_JSON) {
		r.requestParams.body = r.Body
	} else if strings.Contains(hv, HEADER_APPLICATION_XML) {
		r.requestParams.body = r.Body
	}

	r.requestParams.pathParams = mux.Vars(r.Request)
	r.requestParams.getParams = vsTov(r.URL.Query())
	r.requestParams.params = r.requestParams.postParams
	if r.requestParams.params == nil {
		r.requestParams.params = make(map[string]string)
	}

	//合并参数
	for k, v := range r.requestParams.getParams {
		_, ok := r.requestParams.params[k]
		if !ok {
			r.requestParams.params[k] = v
		}
	}

	for k, v := range r.requestParams.pathParams {
		_, ok := r.requestParams.params[k]
		if !ok {
			r.requestParams.params[k] = v
		}
	}

	r.requestParams.state = 1
	r.requestParams.mux.Unlock()
	return nil
}

func vsTov(vs map[string][]string) map[string]string {
	res := make(map[string]string)
	for k, v := range vs {
		res[k] = v[0]
	}
	return res
}
