package csvx_test

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/prongbang/csvx"
	"log"
	"testing"
)

type MyStruct struct {
	Name string `json:"name" db:"name" header:"Name Space" no:"2"`
	ID   int    `json:"id" db:"id" header:"ID" no:"1"`
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

func TestConvertIgnoreDoubleQuote(t *testing.T) {
	// Given
	m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
	expected := `ID,Name Space
1,N1
2,N2`

	// When
	csv := csvx.Convert[MyStruct](m, true)

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

/*
BenchmarkTryConvert
BenchmarkTryConvert-10    	  234032	      5198 ns/op
*/
func BenchmarkTryConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
		expected := `"ID","Name Space"
"1","N1"
"2","N2"`

		// When
		csv := csvx.TryConvert[MyStruct](m)

		// Then
		if csv != expected {
			b.Error("Convert error:", csv)
		}
	}
}

/*
BenchmarkConvertToCSV
BenchmarkConvertToCSV-10    	  584448	      2069 ns/op
*/
func BenchmarkConvertToCSV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
		expected := `ID,Name Space
1,N1
2,N2
`

		// When
		csv := csvx.ConvertToCSV[MyStruct](m)

		// Then
		if csv != expected {
			b.Error("Convert error:", csv, expected)
		}
	}
}

/*
BenchmarkConvertToCSV2
BenchmarkConvertToCSV2-10    	  642216	      1853 ns/op
*/
func BenchmarkConvertToCSV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
		expected := `ID,Name Space
1,N1
2,N2
`

		// When
		csv := csvx.ConvertToCSV2[MyStruct](m)

		// Then
		if csv != expected {
			b.Error("Convert error:", csv, expected)
		}
	}
}

/*
BenchmarkWriter
BenchmarkWriter-10    	 1936059	       581.2 ns/op
*/
func BenchmarkWriter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
		expected := `ID,Name Space
1,N1
2,N2
`

		// When
		var buffer bytes.Buffer
		w := csv.NewWriter(&buffer)

		// Write headers
		headers := []string{"ID", "Name Space"}

		w.Write(headers)
		for _, record := range m {
			row := []string{fmt.Sprint(record.ID), record.Name}
			if err := w.Write(row); err != nil {
				log.Fatalln("error writing record to file", err)
			}
		}
		w.Flush()
		csvs := buffer.String()

		// Then
		if csvs != expected {
			b.Error("Convert error:", csvs)
		}
	}
}
