package csvx_test

import (
	"fmt"
	"testing"

	"github.com/prongbang/csvx"
)

type MyStruct struct {
	Name  string `json:"name" db:"name" header:"Name Space" no:"2"`
	ID    int    `json:"id" db:"id" header:"ID" no:"1"`
	Other string
}

type MyStructPointer struct {
	Name    string  `json:"name" db:"name" header:"Name Space" no:"2"`
	ID      int     `json:"id" db:"id" header:"ID" no:"1"`
	Address *string `json:"address" db:"address" header:"Address" no:"4" default:"N/A"`
	Phone   *string `json:"phone" db:"phone" header:"Phone" no:"3" default:"NULL"`
	Email   *string `json:"email" db:"email" header:"Email" no:"5"`
	Age     *string `json:"age" db:"age" header:"Age" no:"6" default:"999"`
	Other   string
}

func TestConvert(t *testing.T) {
	// Given
	m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
	expected := csvx.Utf8BOM + `"ID","Name Space"
"1","N1"
"2","N2"`

	// When
	result := csvx.Convert[MyStruct](m)

	// Then
	if result != expected {
		t.Error("Convert error:", result)
	}
}

func TestTryConvertPointer(t *testing.T) {
	// Given
	age := "100"
	phone := "0876"
	m := []MyStructPointer{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2", Age: &age, Phone: &phone}}
	expected1 := csvx.Utf8BOM + `"ID","Name Space","Phone","Address","Email","Age"
"1","N1","NULL","N/A","","999"
"2","N2","0876","N/A","","100"`
	expected2 := csvx.Utf8BOM + `ID,Name Space,Phone,Address,Email,Age
"1","N1","NULL","N/A","","999"
"2","N2","0876","N/A","","100"
`

	// When
	result1 := csvx.Convert(m)
	result2 := csvx.TryConvert(m)

	// Then
	if result1 != expected1 {
		t.Error("Result1 error:\nexpected:", expected1, "\nactual:", result1)
	}
	if result2 != expected2 {
		t.Error("Result2 error:\nexpected:", expected2, "\nactual:", result2)
	}
}

func TestConvertIgnoreDoubleQuote(t *testing.T) {
	// Given
	m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
	expected := csvx.Utf8BOM + `ID,Name Space
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
		expected := csvx.Utf8BOM + `"ID","Name Space"
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
		expected := csvx.Utf8BOM + `ID,Name Space
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
		expected := csvx.Utf8BOM + `ID,Name Space
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
