package fastgo

import (
	"fastgo/base_type"
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
	name     string
	kind     reflect.Kind
	defValue string
	tag      reflect.StructTag
}

func (c *reflectCache) get(t reflect.Type) *reflectValue {

	c.mux.RLock()
	rv, ok := c.cache[t]
	if !ok {
		rv = c.set(t)
	}
	c.mux.RUnlock()
	return rv
}

func (c *reflectCache) set(t reflect.Type) *reflectValue {
	elem := t.Elem()

	rv := new(reflectValue)

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
	return rv
}

func Bind(p any, values map[string]string) error {
	var paramValue reflect.Value

	paramValue = reflect.ValueOf(p)
	if paramValue.Kind() == reflect.Pointer {
		paramValue = paramValue.Elem()
	}

	rc := rCache.get(reflect.TypeOf(p))

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
				err := Bind(fieldValue.Addr().Interface(), values)
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
