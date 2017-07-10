package util

import (
	"reflect"
	"strconv"
	"errors"
	"time"
)

func Apply(s reflect.Kind, val string) (interface{}, error) {
	switch s {
	case reflect.Int:
		return strconv.ParseInt(val, 10, 0)
	case reflect.Int8:
		return strconv.ParseInt(val, 10, 8)
	case reflect.Int16:
		return strconv.ParseInt(val, 10, 16)
	case reflect.Int32:
		return strconv.ParseInt(val, 10, 32)
	case reflect.Int64:
		return strconv.ParseInt(val, 10, 64)
	case reflect.Uint:
		return strconv.ParseUint(val, 10, 0)
	case reflect.Uint8:
		return strconv.ParseUint(val, 10, 8)
	case reflect.Uint16:
		return strconv.ParseUint(val, 10, 16)
	case reflect.Uint32:
		return strconv.ParseUint(val, 10, 32)
	case reflect.Uint64:
		return strconv.ParseUint(val, 10, 64)
	case reflect.Bool:
		return strconv.ParseBool(val)
	case reflect.Float32:
		return strconv.ParseFloat(val, 32)
	case reflect.Float64:
		return strconv.ParseFloat(val, 64)
	case reflect.String:
		return val, nil
	case reflect.Struct:
		// TODO Struct Kind but not time.Time
		return time.Parse(time.RFC3339, val)
	default:
		return nil, errors.New("unknown type")
	}
}