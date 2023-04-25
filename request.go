package fastgo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
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

func (r *Request) Params(param any) {

	if r.Method == http.MethodGet {
		r.bind(param, r.URL.Query())
		return
	}
	if r.Method == http.MethodPost {
		hv := r.Header.Get("Content-Type")
		if strings.Contains(hv, MULTIPART_FORM_DATA) {
			err := r.ParseMultipartForm(defaultMaxMemory)
			if err != nil {
				return
			}
		}

		if strings.Contains(hv, X_WWW_FORM_URLENCODED) {
			err := r.ParseForm()
			if err != nil {
				return
			}
		}

		if strings.Contains(hv, APPLICATION_JSON) {
			err := json.NewDecoder(r.Body).Decode(param)
			if err != nil {
				return
			}
		}

		r.bind(param, r.Form)
		return
	}
	return
}

func (r *Request) bind(param any, values url.Values) {
	paramValue := reflect.ValueOf(param).Elem()
	for i := 0; i < paramValue.Type().NumField(); i++ {

		field := paramValue.Type().Field(i)
		vs, ok := values[field.Name]
		if !ok {
			continue
		}

		if vs == nil {
			continue
		}

		fieldValue := paramValue.FieldByName(field.Name)
		if !fieldValue.CanSet() {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Struct:
		case reflect.Map:
		case reflect.Slice:
		case reflect.String:
			fieldValue.SetString(vs[0])
		case reflect.Int64:
			i, _ := strconv.ParseInt(vs[0], 10, 64)
			fieldValue.SetInt(i)
		}
	}
}

//func parseType(src reflect.Value, tar reflect.Value) {
//
//	p := src.UnsafePointer()
//	switch src.Kind() {
//	case Int:
//		return int64(*(*int)(p))
//	case Int8:
//		return int64(*(*int8)(p))
//	case Int16:
//		return int64(*(*int16)(p))
//	case Int32:
//		return int64(*(*int32)(p))
//	case Int64:
//	}
//}
