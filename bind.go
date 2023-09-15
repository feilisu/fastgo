package fastgo

import (
	"fastgo/base_type"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var (
	rCache = new(reflectCache)
)

type reflectCache struct {
	cache map[reflect.Type]*reflectValue
	mux   sync.RWMutex
}

type reflectValue struct {
	fields []*reflectField
}

type reflectField struct {
	name     string
	kind     reflect.Kind
	defValue string
	tag      reflect.StructTag
}

func (c *reflectCache) get(t reflect.Type) (rv *reflectValue, err error) {

	c.mux.RLock()
	rv, ok := c.cache[t]
	if !ok {
		rv, err = c.set(t)
	}
	c.mux.RUnlock()
	return
}

func (c *reflectCache) set(t reflect.Type) (rv *reflectValue, err error) {
	if t.Kind() != reflect.Struct {
		return nil, NewError("类型必须是struct")
	}

	elem := t.Elem()
	rv = new(reflectValue)

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		rf := &reflectField{
			name:     field.Name,
			kind:     field.Type.Kind(),
			tag:      field.Tag,
			defValue: field.Tag.Get("default"),
		}
		rv.fields = append(rv.fields, rf)
	}

	if rCache.cache == nil {
		rCache.cache = make(map[reflect.Type]*reflectValue)
	}

	rCache.cache[t] = rv
	return
}

type Binding struct {
}

func NewBinding() *Binding {
	return new(Binding)
}

func (bi *Binding) Bind(p any, values map[string]string) error {
	var paramValue reflect.Value

	paramValue = reflect.ValueOf(p)
	if paramValue.Kind() == reflect.Pointer {
		paramValue = paramValue.Elem()
	}

	rc, err := rCache.get(reflect.TypeOf(p))
	if err != nil {
		return err
	}

	for _, rv := range rc.fields {

		field, b := paramValue.Type().FieldByName(rv.name)
		if !b {
			continue
		}

		vs, ok := values[base_type.FirstToLow(field.Name)]
		if !ok {
			//是否设置了默认值
			if rv.defValue == "" {
				continue
			}
			vs = rv.defValue
		}

		fieldValue := paramValue.FieldByName(field.Name)
		if !fieldValue.CanSet() {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			if fieldValue.CanAddr() {
				err := bi.Bind(fieldValue.Addr().Interface(), values)
				if err != nil {
					return err
				}
			}
		//case reflect.Map:
		//case reflect.Slice:
		//	fieldValue.Type()
		//	if fieldValue.CanAddr() {
		//		err := json.Unmarshal([]byte(vs), fieldValue.Pointer())
		//		if err != nil {
		//			return err
		//		}
		//	}
		case reflect.Bool:
			parseBool, _ := strconv.ParseBool(vs)
			fieldValue.SetBool(parseBool)
		case reflect.String:
			fieldValue.SetString(vs)
		case reflect.Int64:
			_ = stringToInt(vs, 64, fieldValue)
		case reflect.Int32:
			_ = stringToInt(vs, 32, fieldValue)
		case reflect.Int16:
			_ = stringToInt(vs, 16, fieldValue)
		case reflect.Int8:
			_ = stringToInt(vs, 8, fieldValue)
		case reflect.Int:
			_ = stringToInt(vs, 0, fieldValue)
		case reflect.Uint64:
			_ = stringToUint(vs, 64, fieldValue)
		case reflect.Uint32:
			_ = stringToUint(vs, 32, fieldValue)
		case reflect.Uint16:
			_ = stringToUint(vs, 16, fieldValue)
		case reflect.Uint8:
			_ = stringToUint(vs, 8, fieldValue)
		case reflect.Uint:
			_ = stringToUint(vs, 0, fieldValue)
		case reflect.Float64:
			_ = stringToFloat(vs, 64, fieldValue)
		case reflect.Float32:
			_ = stringToFloat(vs, 32, fieldValue)
		}
	}
	return nil
}

func stringToInt(str string, bitSize int, value reflect.Value) error {
	i, err := strconv.ParseInt(str, 10, bitSize)
	if err != nil {
		return err
	}
	value.SetInt(i)
	return nil
}

func stringToUint(str string, bitSize int, value reflect.Value) error {
	i, err := strconv.ParseUint(str, 10, bitSize)
	if err != nil {
		return err
	}
	value.SetUint(i)
	return nil
}

func stringToFloat(str string, bitSize int, value reflect.Value) error {
	i, err := strconv.ParseFloat(str, bitSize)
	if err != nil {
		return err
	}
	value.SetFloat(i)
	return nil
}

type Validator struct {
}

func NewValidate() *Validator {
	return new(Validator)
}

func (v *Validator) validate(p any) error {

	rc, err := rCache.get(reflect.TypeOf(p))
	if err != nil {
		return err
	}

	var paramValue reflect.Value

	paramValue = reflect.ValueOf(p)
	if paramValue.Kind() == reflect.Pointer {
		paramValue = paramValue.Elem()
	}

	for _, rv := range rc.fields {

		fieldValue := paramValue.FieldByName(rv.name)
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i := _int64(fieldValue.Int())
			err = i.validate(rv.tag.Get("validate"), rv.name, "")
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i := _uint64(fieldValue.Int())
			err = i.validate(rv.tag.Get("validate"), rv.name, "")
		case reflect.Float32, reflect.Float64:
			i := _float64(fieldValue.Int())
			err = i.validate(rv.tag.Get("validate"), rv.name, "")
		case reflect.Bool:
		case reflect.String:
			i := _string(fieldValue.String())
			err = i.validate(rv.tag.Get("validate"), rv.name, "")
		case reflect.Map:
		case reflect.Slice:
		case reflect.Struct:
		}

		if err != nil {
			return nil
		}
	}
	return nil
}

type _int64 int64
type _uint64 uint64
type _float64 float64
type _string string

func parserTag(tag string) [][]string {
	if len(tag) <= 0 {
		return nil
	}

	var strss [][]string
	tags := strings.Split(tag, ",")
	for _, s := range tags {
		ss := strings.Split(s, "=")
		strss = append(strss, ss)
	}
	return strss
}

func (i *_int64) validate(tag string, fieldName string, msg string) error {

	strss := parserTag(tag)
	fieldValue := int64(*i)

	for i := 0; i < len(strss); i++ {
		var rn string
		var rv int64
		r := strss[i]
		rn = r[0]
		if len(r) > 1 {
			rv, _ = strconv.ParseInt(r[1], 10, 64)
		}

		var errorMsg *Error

		switch rn {
		case "gt":
			if fieldValue <= rv {
				errorMsg = NewError(fieldName + "应该大于" + r[1])
			}
		case "gte":
			if fieldValue < rv {
				errorMsg = NewError(fieldName + "应该大于等于" + r[1])
			}
		case "lt":
			if fieldValue >= rv {
				errorMsg = NewError(fieldName + "应该小于" + r[1])
			}
		case "lte":
			if fieldValue > rv {
				errorMsg = NewError(fieldName + "应该小于等于" + r[1])
			}

		}

		if errorMsg != nil {

			if len(msg) > 0 {
				errorMsg.SetMsg(msg)
			}
			return errorMsg
		}

	}

	return nil
}

func (i *_uint64) validate(tag string, fieldName string, msg string) error {

	strss := parserTag(tag)
	fieldValue := uint64(*i)

	for i := 0; i < len(strss); i++ {
		var rn string
		var rv uint64
		r := strss[i]
		rn = r[0]
		if len(r) > 1 {
			rv, _ = strconv.ParseUint(r[1], 10, 64)
		}

		var errorMsg *Error

		switch rn {
		case "gt":
			if fieldValue <= rv {
				errorMsg = NewError(fieldName + "应该大于" + r[1])
			}
		case "gte":
			if fieldValue < rv {
				errorMsg = NewError(fieldName + "应该大于等于" + r[1])
			}
		case "lt":
			if fieldValue >= rv {
				errorMsg = NewError(fieldName + "应该小于" + r[1])
			}
		case "lte":
			if fieldValue > rv {
				errorMsg = NewError(fieldName + "应该小于等于" + r[1])
			}

		}

		if errorMsg != nil {

			if len(msg) > 0 {
				errorMsg.SetMsg(msg)
			}
			return errorMsg
		}

	}

	return nil
}

func (i *_float64) validate(tag string, fieldName string, msg string) error {

	strss := parserTag(tag)
	fieldValue := float64(*i)

	for i := 0; i < len(strss); i++ {
		var rn string
		var rv float64
		r := strss[i]
		rn = r[0]
		if len(r) > 1 {
			rv, _ = strconv.ParseFloat(r[1], 64)
		}

		var errorMsg *Error

		switch rn {
		case "gt":
			if fieldValue <= rv {
				errorMsg = NewError(fieldName + "应该大于" + r[1])
			}
		case "gte":
			if fieldValue < rv {
				errorMsg = NewError(fieldName + "应该大于等于" + r[1])
			}
		case "lt":
			if fieldValue >= rv {
				errorMsg = NewError(fieldName + "应该小于" + r[1])
			}
		case "lte":
			if fieldValue > rv {
				errorMsg = NewError(fieldName + "应该小于等于" + r[1])
			}

		}

		if errorMsg != nil {

			if len(msg) > 0 {
				errorMsg.SetMsg(msg)
			}
			return errorMsg
		}

	}

	return nil
}

func (i *_string) validate(tag string, fieldName string, msg string) error {

	strss := parserTag(tag)
	fieldValue := string(*i)
	length := int64(len(fieldValue))

	for i := 0; i < len(strss); i++ {
		var rn string
		var rv int64
		r := strss[i]
		rn = r[0]
		if len(r) > 1 {
			rv, _ = strconv.ParseInt(r[1], 10, 64)
		}

		var errorMsg *Error

		switch rn {
		case "gt":
			if length <= rv {
				errorMsg = NewError(fieldName + "长度应该大于" + r[1])
			}
		case "gte":
			if length < rv {
				errorMsg = NewError(fieldName + "长度应该大于等于" + r[1])
			}
		case "lt":
			if length >= rv {
				errorMsg = NewError(fieldName + "长度应该小于" + r[1])
			}
		case "lte":
			if length > rv {
				errorMsg = NewError(fieldName + "长度应该小于等于" + r[1])
			}

		}

		if errorMsg != nil {

			if len(msg) > 0 {
				errorMsg.SetMsg(msg)
			}
			return errorMsg
		}

	}

	return nil
}
