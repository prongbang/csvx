package csvx_test

import (
	"fmt"
	"github.com/prongbang/csvx"
	"testing"
)

type MyStruct struct {
	Name  string `json:"name" db:"name" header:"Name Space" no:"2"`
	ID    int    `json:"id" db:"id" header:"ID" no:"1"`
	Other string
}

func TestConvert(t *testing.T) {
	// Given
	m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
	expected := `"ID","Name Space"
"1","N1"
"2","N2"`

	// When
	result := csvx.Convert[MyStruct](m)

	// Then
	if result != expected {
		t.Error("Convert error:", result)
	}
}

func TestConvertIgnoreDoubleQuote(t *testing.T) {
	// Given
	m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
	expected := `ID,Name Space
1,N1
2,N2`

	// When
	result := csvx.Convert[MyStruct](m, true)

	// Then
	if result != expected {
		t.Error("Convert error:", result)
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
		result := csvx.Convert[MyStruct](m)

		// Then
		if result != expected {
			b.Error("Convert error:", result)
		}
	}
}

func BenchmarkManualConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
		expected := `ID,Name Space
1,N1
2,N2
`

		// When
		result := csvx.ManualConvert[MyStruct](m,
			[]string{"ID", "Name Space"},
			func(data MyStruct) []string {
				return []string{
					fmt.Sprintf("%d", data.ID),
					data.Name,
				}
			},
		)

		// Then
		if result != expected {
			b.Error("Convert error:", result)
		}
	}
}

func BenchmarkTryConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
		expected := `ID,Name Space
"1","N1"
"2","N2"
`

		// When
		result := csvx.TryConvert(m)

		// Then
		if result != expected {
			b.Error("Convert error:", result)
		}
	}
}
