package csvx

import (
	"reflect"
	"strconv"
)

func IsFloat(t reflect.Type) bool {
	if t.Kind() == reflect.Float32 || t == reflect.TypeOf(float64(0)) {
		return true
	}
	return false
}

func F64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
