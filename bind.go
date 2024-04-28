package fastgo

import (
	"encoding/json"
	"fastgo/internal/util"
	"fmt"
	"reflect"
	"strconv"
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
	name string
	kind reflect.Kind
	tag  reflect.StructTag
}

func (c *reflectCache) get(t reflect.Type) (rv *reflectValue, err error) {

	c.mux.Lock()
	rv, ok := c.cache[t]
	if !ok {
		rv, err = c.set(t)
	}
	c.mux.Unlock()
	return
}

func (c *reflectCache) set(t reflect.Type) (rv *reflectValue, err error) {

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, NewError("类型必须是struct")
	}

	rv = new(reflectValue)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		rf := &reflectField{
			name: field.Name,
			kind: field.Type.Kind(),
			tag:  field.Tag,
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
	var paramType reflect.Type

	paramValue = reflect.ValueOf(p)
	if paramValue.Kind() == reflect.Pointer {
		paramValue = paramValue.Elem()
	}
	paramType = paramValue.Type()

	rc, err := rCache.get(paramType)
	if err != nil {
		return err
	}

	for _, rv := range rc.fields {

		field, b := paramType.FieldByName(rv.name)
		if !b {
			continue
		}

		vs, ok := values[util.FirstToLow(field.Name)]
		if !ok {
			continue
		}

		fieldValue := paramValue.FieldByName(field.Name)
		if !fieldValue.CanSet() {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			if fieldValue.CanAddr() {
				err := json.Unmarshal([]byte(vs), fieldValue.Addr().Interface())
				if err != nil {
					return err
				}
			}
		case reflect.Map:
			if fieldValue.CanAddr() {
				err := json.Unmarshal([]byte(vs), fieldValue.Addr().Interface())
				if err != nil {
					return err
				}
			}
		case reflect.Slice:
			if fieldValue.CanAddr() {
				err := json.Unmarshal([]byte(vs), fieldValue.Addr().Interface())
				if err != nil {
					return err
				}
			}
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

	if valid, ok := p.(Validator); ok {
		err := valid.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (bi *Binding) BindRequestParams(p any, values map[string]any) error {
	if values == nil {
		return NewError("bind values is nil ")
	}

	var paramValue reflect.Value
	var paramType reflect.Type

	paramValue = reflect.ValueOf(p)
	if paramValue.Kind() == reflect.Pointer {
		paramValue = paramValue.Elem()
	}
	paramType = paramValue.Type()

	rc, err := rCache.get(paramType)
	if err != nil {
		return err
	}

	for _, rv := range rc.fields {

		field, b := paramType.FieldByName(rv.name)
		if !b {
			continue
		}

		vs, ok := values[util.FirstToLow(field.Name)]
		if !ok {
			continue
		}

		fieldValue := paramValue.FieldByName(field.Name)
		if !fieldValue.CanSet() {
			continue
		}

		switch field.Type.Kind() {
		//case reflect.Struct:
		//	if fieldValue.CanAddr() {
		//		vsm := anyToMap(vs)
		//		if vsm != nil {
		//			err := bi.BindRequestParams(fieldValue.Addr().Interface(), vsm)
		//			if err != nil {
		//				return err
		//			}
		//		}
		//	}
		//case reflect.Map:
		//	if fieldValue.CanAddr() {
		//		vsm := anyToMap(vs)
		//		if vsm != nil {
		//			err := bi.BindRequestParams(fieldValue.Addr().Interface(), vsm)
		//			if err != nil {
		//				return err
		//			}
		//		}
		//	}
		////case reflect.Slice:
		////	fieldValue.Type()
		////	if fieldValue.CanAddr() {
		////		err := json.Unmarshal([]byte(vs), fieldValue.Pointer())
		////		if err != nil {
		////			return err
		////		}
		////	}
		////	fieldValue.set
		case reflect.Bool:
			b, err := anyToBool(vs)
			if err != nil {
				return err
			}
			fieldValue.SetBool(b)
		case reflect.String:
			s, err := anyToString(vs)
			if err != nil {
				return err
			}
			fieldValue.SetString(s)
		case reflect.Int64:
			i, err := anyToInt(vs, 64)
			if err != nil {
				return err
			}
			fieldValue.SetInt(i)
		case reflect.Int32:
			i, err := anyToInt(vs, 32)
			if err != nil {
				return err
			}
			fieldValue.SetInt(i)
		case reflect.Int16:
			i, err := anyToInt(vs, 16)
			if err != nil {
				return err
			}
			fieldValue.SetInt(i)
		case reflect.Int8:
			i, err := anyToInt(vs, 8)
			if err != nil {
				return err
			}
			fieldValue.SetInt(i)
		case reflect.Int:
			i, err := anyToInt(vs, 0)
			if err != nil {
				return err
			}
			fieldValue.SetInt(i)
		case reflect.Uint64:
			i, err := anyToUint(vs)
			if err != nil {
				return err
			}
			fieldValue.SetUint(i)
		case reflect.Uint32:
			i, err := anyToUint(vs)
			if err != nil {
				return err
			}
			fieldValue.SetUint(i)
		case reflect.Uint16:
			i, err := anyToUint(vs)
			if err != nil {
				return err
			}
			fieldValue.SetUint(i)
		case reflect.Uint8:
			i, err := anyToUint(vs)
			if err != nil {
				return err
			}
			fieldValue.SetUint(i)
		case reflect.Uint:
			i, err := anyToUint(vs)
			if err != nil {
				return err
			}
			fieldValue.SetUint(i)
		case reflect.Float64:
			f, err := anyToFloat(vs)
			if err != nil {
				return err
			}
			fieldValue.SetFloat(f)
		case reflect.Float32:
			f, err := anyToFloat(vs)
			if err != nil {
				return err
			}
			fieldValue.SetFloat(f)
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

func anyToMap(a any) map[string]any {
	m, ok := a.(map[string]any)
	if ok {
		return m
	}
	return nil
}

func anyToInt(a any, bitSize int) (int64, error) {
	k := reflect.TypeOf(a).Kind()
	switch k {
	case reflect.String:
		s := a.(string)
		i, err := strconv.ParseInt(s, 10, bitSize)
		if err != nil {
			return 0, err
		}
		return i, err
	case reflect.Int64:
		u := a.(int64)
		return u, nil
	case reflect.Int32:
		if bitSize < 32 {
			return 0, NewError(fmt.Sprintf("%T can not conversion type to int%d", a, bitSize))
		}
		u := a.(int32)
		return int64(u), nil
	case reflect.Int16:
		if bitSize < 16 {
			return 0, NewError(fmt.Sprintf("%T can not conversion type to int%d", a, bitSize))
		}
		u := a.(int16)
		return int64(u), nil
	case reflect.Int8:
		if bitSize < 8 {
			return 0, NewError(fmt.Sprintf("%T can not conversion type to int%d", a, bitSize))
		}
		u := a.(int8)
		return int64(u), nil
	case reflect.Int:
		sysSize := strconv.IntSize
		if bitSize < sysSize {
			return 0, NewError(fmt.Sprintf("%d system %T can not conversion type to int%d", sysSize, a, bitSize))
		}
		u := a.(uint)
		return int64(u), nil
	default:
		return 0, NewError(fmt.Sprintf("%T can not conversion type to int", a))
	}
}

func anyToFloat(a any) (res float64, err error) {
	k := reflect.TypeOf(a).Kind()
	switch k {
	case reflect.String:
		var f float64
		s := a.(string)
		f, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return
		}
		res = f
		return
	default:
		if reflect.TypeOf(res).Kind() == k {
			res = a.(float64)
			return
		}
		err = NewError(fmt.Sprintf("%T can not conversion type to float", a))
		return
	}
}

func anyToUint(a any) (res uint64, err error) {

	k := reflect.TypeOf(a).Kind()
	switch k {
	case reflect.String:
		var i uint64
		s := a.(string)
		i, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			return
		}
		res = i
		return
	default:
		if reflect.TypeOf(res).Kind() == k {
			res = a.(uint64)
			return
		}
		err = NewError(fmt.Sprintf("%T can not conversion type to uint", a))
		return
	}
}

func anyToString(a any) (string, error) {
	k := reflect.TypeOf(a).Kind()
	switch k {
	case reflect.String:
		s := a.(string)
		return s, nil
	default:
		return "", NewError(fmt.Sprintf("%T can not conversion type to string", a))
	}
}

func anyToBool(a any) (bool, error) {
	k := reflect.TypeOf(a).Kind()
	if k == reflect.String {
		s := a.(string)
		b, err := strconv.ParseBool(s)
		if err != nil {
			return b, nil
		}
		return b, nil
	} else if k == reflect.Bool {
		b := a.(bool)
		return b, nil
	} else {
		return false, NewError(fmt.Sprintf("%T can not conversion type to bool", a))
	}
}

// Validator 验证器
type Validator interface {
	Validate() error
}
