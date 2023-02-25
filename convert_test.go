package csvx_test

import (
	"github.com/prongbang/csvx"
	"testing"
)

type MyStruct struct {
	Name string `json:"name" db:"name" field:"Name Space" index:"2"`
	ID   int    `json:"id" db:"id" field:"ID" index:"1"`
}

func TestConvert(t *testing.T) {
	// Given
	m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
	expected := `"ID","Name Space"
"1","N1"
"2","N2"`

	// When
	csv := csvx.Convert[MyStruct](m)

	// Then
	if csv != expected {
		t.Error("Convert error:", csv)
	}
}

func BenchmarkConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
		expected := `"ID","Name Space"
"1","N1"
"2","N2"`

		// When
		csv := csvx.Convert[MyStruct](m)

		// Then
		if csv != expected {
			b.Error("Convert error:", csv)
		}
	}
}
