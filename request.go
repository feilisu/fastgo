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

func (r *Request) Params(param any) error {

	if r.Method == http.MethodGet {
		return r.bind(param, r.URL.Query())
	}
	if r.Method == http.MethodPost {
		hv := r.Header.Get("Content-Type")
		if strings.Contains(hv, MULTIPART_FORM_DATA) {
			if err := r.ParseMultipartForm(defaultMaxMemory); err != nil {
				return err
			}
			return r.bind(param, r.MultipartForm.Value)
		}

		if strings.Contains(hv, X_WWW_FORM_URLENCODED) {
			if err := r.ParseForm(); err != nil {
				return err
			}
			return r.bind(param, r.Form)
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

func (r *Request) bind(param any, values url.Values) error {
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
			reflect.MakeMap()
			field.Type
			decodeJson(vs[0])
		case reflect.Map:
		case reflect.Slice:
		case reflect.Bool:
			parseBool, _ := strconv.ParseBool(vs[0])
			fieldValue.SetBool(parseBool)
		case reflect.String:
			fieldValue.SetString(vs[0])
		case reflect.Int64:
			_ = stringToInt(vs[0], 64, fieldValue)
		case reflect.Int32:
			_ = stringToInt(vs[0], 32, fieldValue)
		case reflect.Int16:
			_ = stringToInt(vs[0], 16, fieldValue)
		case reflect.Int8:
			_ = stringToInt(vs[0], 8, fieldValue)
		case reflect.Int:
			_ = stringToInt(vs[0], 0, fieldValue)
		case reflect.Uint64:
			_ = stringToUint(vs[0], 64, fieldValue)
		case reflect.Uint32:
			_ = stringToUint(vs[0], 32, fieldValue)
		case reflect.Uint16:
			_ = stringToUint(vs[0], 16, fieldValue)
		case reflect.Uint8:
			_ = stringToUint(vs[0], 8, fieldValue)
		case reflect.Uint:
			_ = stringToUint(vs[0], 0, fieldValue)
		case reflect.Float64:
			_ = stringToFloat(vs[0], 64, fieldValue)
		case reflect.Float32:
			_ = stringToFloat(vs[0], 32, fieldValue)
		}
	}
	return nil
}

func stringToInt(str string, bitSize int, value reflect.Value) error {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	value.SetInt(i)
	return nil
}

func stringToUint(str string, bitSize int, value reflect.Value) error {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}
	value.SetUint(i)
	return nil
}

func stringToFloat(str string, bitSize int, value reflect.Value) error {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	value.SetFloat(i)
	return nil
}

func decodeJson(bytes string, param any) error {
	err := json.Unmarshal([]byte(bytes), param)
	if err != nil {
		return err
	}
	return nil
}
