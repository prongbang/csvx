package csvx_test

import (
	"github.com/prongbang/csvx"
	"reflect"
	"testing"
)

func TestIsFloat(t *testing.T) {
	ty := reflect.TypeOf(float64(0))
	if !csvx.IsFloat(ty) {
		t.Error("Is Not Float")
	}
	if !csvx.IsFloat(reflect.TypeOf(float32(0))) {
		t.Error("Float32 is not float")
	}
	if csvx.IsFloat(reflect.TypeOf(1)) {
		t.Error("Int is float")
	}
}

func TestIsPointer(t *testing.T) {
	value := ""
	if !csvx.IsPointer(reflect.TypeOf(&value)) {
		t.Error("Is not pointer")
	}
	if csvx.IsPointer(reflect.TypeOf(value)) {
		t.Error("String is pointer")
	}
}

func TestF64ToString(t *testing.T) {
	f := 6.4
	if csvx.F64ToString(f) != "6.4" {
		t.Error("Actual is not eq", f)
	}
}

func TestRemoveDoubleQuote(t *testing.T) {
	text := "\ufeff\"hello\" world\""
	expect := "hello\" world"
	if csvx.RemoveDoubleQuote(text) != expect {
		t.Error("Actual is not eq", expect)
	}
}

func TestRemoveDoubleQuoteWithoutWrappingQuotes(t *testing.T) {
	text := "\ufeffhello"
	expect := "hello"
	if csvx.RemoveDoubleQuote(text) != expect {
		t.Error("Actual is not eq", expect)
	}
}
