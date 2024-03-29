package csvx

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// RemoveDoubleQuote remove double quote (") and clear unicode from text
func RemoveDoubleQuote(text string) string {
	text = ClearUnicode(text)
	first := strings.Index(text, "\"")
	last := strings.LastIndex(text, "\"")
	if first == 0 && last == (len(text)-1) {
		text = text[1:last]
	}
	return text
}

// ClearUnicode clear unicode from text
func ClearUnicode(text string) string {
	regex := regexp.MustCompile("^\ufeff")
	result := regex.ReplaceAllString(text, "")
	return strings.TrimSpace(result)
}

// IsFloat returns true if the given reflect.Type is a float32 or float64 type, and false otherwise.
// This function can be used to check whether a given type is a floating-point type, which may be useful
// for type assertions and other operations that require type checking. If the given type is not a valid
// float type, this function will return false.
func IsFloat(t reflect.Type) bool {
	if t.Kind() == reflect.Float32 || t == reflect.TypeOf(float64(0)) {
		return true
	}
	return false
}

func IsPointer(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer || t.Kind() == reflect.Ptr
}

// F64ToString converts the given float64 value to a string representation.
// The resulting string will be formatted as a decimal number with up to 10 decimal places.
// This function can be used to convert floating-point values to string values, which may be useful
// for printing or other output operations. If the given value is NaN or infinite, the resulting string will reflect this.
// If the given value is not representable as a finite decimal number, this function may return an inaccurate or nonsensical result.
func F64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
