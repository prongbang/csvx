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
