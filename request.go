package fastgo

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type Request struct {
	*http.Request
}

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

func (r *Request) params() url.Values {

	if r.Method == http.MethodGet {
		return r.URL.Query()
	}
	if r.Method == http.MethodPost {
		hv := r.Header.Get("Content-Type")
		if strings.Contains(hv, "multipart/form-data") {
			err := r.ParseMultipartForm(30 << 20)
			if err != nil {
				return nil
			}
		}

		if strings.Contains(hv, "application/x-www-form-urlencoded") {
			err := r.ParseForm()
			if err != nil {
				return nil
			}
		}

		//if strings.Contains(hv, "application/json") {
		//
		//}
		return r.Form
	}
	return nil
}

func (r *Request) PostJsonParams(par any) any {
	var pp any

	err := json.NewDecoder(r.Body).Decode(&pp)
	if err != nil {
		return nil
	}
	tar := reflect.ValueOf(par).Elem()
	src := reflect.ValueOf(pp)
	if tar.Kind() == reflect.Struct {
		bindStruct(src, tar)
	}
	log.Println(src.Kind())
	log.Println(tar.Kind())
	return pp
}

func bindStruct(src reflect.Value, tar reflect.Value) {
	if src.Kind() != reflect.Map {
		return
	}

	mkv := make(map[string]reflect.Value)
	mapIterator := src.MapRange()
	for mapIterator.Next() {
		mkv[mapIterator.Key().Type().Name()] = mapIterator.Value()
	}

	for i := 0; i < tar.Type().NumField(); i++ {

		field := tar.Type().Field(i)
		v, ok := mkv[field.Name]
		if !ok {
			continue
		}

		log.Println()
		switch field.Type.Kind() {
		case reflect.Struct:
			bindStruct(v, field)
		case reflect.Map:
			bindMap(v, field.Type)
		case reflect.Slice:
			bindSplice(v, field.Type)
		case reflect.String:

		case reflect.Int64:

		}
		//src.FieldByName()
	}

}
func bindMap(src reflect.Value, tar reflect.Value) {

}

func bindSplice(src reflect.Value, tar reflect.Value) {

}
