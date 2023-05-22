package fastgo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	*http.Request
}

const (
	APPLICATION_JSON      = "application/json"
	TEXT_HTML             = "text/html"
	APPLICATION_XML       = "application/xml"
	TEXT_XML              = "text/xml"
	TEXT_PLAIN            = "text/plain"
	X_WWW_FORM_URLENCODED = "application/x-www-form-urlencoded"
	MULTIPART_FORM_DATA   = "multipart/form-data"
)

const defaultMaxMemory int64 = 30 << 20

// getString
func (r *Request) getString(key string) string {
	return key
}

func (r *Request) GetJsonBody() string {
	//r.URL.RawQuery
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	return string(bytes)
}

//func (r *Request) Demo() {
//	log.Println(r.postJsonParams())
//}

func (r *Request) Params(param any) error {

	if r.Method == http.MethodGet {
		return NewBinding().Bind(param, convertValues(r.URL.Query()))
	}
	if r.Method == http.MethodPost {
		hv := r.Header.Get("Content-Type")
		if strings.Contains(hv, MULTIPART_FORM_DATA) {
			if err := r.ParseMultipartForm(defaultMaxMemory); err != nil {
				return err
			}
			return NewBinding().Bind(param, convertValues(r.MultipartForm.Value))
		}

		if strings.Contains(hv, X_WWW_FORM_URLENCODED) {
			if err := r.ParseForm(); err != nil {
				return err
			}
			return NewBinding().Bind(param, convertValues(r.Form))
		}

		if strings.Contains(hv, APPLICATION_JSON) {
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(param); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func convertValues(values url.Values) map[string]string {
	m := make(map[string]string, len(values))
	for k, v := range values {
		m[k] = v[0]
	}
	return m
}
