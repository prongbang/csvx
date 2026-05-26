package csvx_test

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/prongbang/csvx"
)

type Struct struct {
	ID   string `header:"ID"`
	Name string `header:"Name Space"`
}

type StructType struct {
	ID   int64   `header:"ID" no:"1"`
	Name string  `header:"Name Space" no:"2"`
	Age  float64 `header:"Age" no:"3"`
}

type StructPointerType struct {
	ID   *int64   `header:"ID" no:"1"`
	Name *string  `header:"Name Space" no:"2"`
	Age  *float64 `header:"Age" no:"3"`
}

type StructAllTypes struct {
	ID     int64   `header:"ID"`
	Count  uint64  `header:"Count"`
	Age    float64 `header:"Age"`
	Active bool    `header:"Active"`
}

type StructAllPointerTypes struct {
	ID     *int64   `header:"ID"`
	Count  *uint64  `header:"Count"`
	Age    *float64 `header:"Age"`
	Active *bool    `header:"Active"`
}

type StructWithStructField struct {
	Value struct{} `header:"Value"`
}

type StructWithPointerStructField struct {
	Value *struct{} `header:"Value"`
}

func TestParserString(t *testing.T) {
	// Given
	rows := [][]string{
		{"\ufeffID", "Name Space"},
		{"1", "Name1"},
		{"2", "Name2"},
		{"3", "Name3"},
		{"4", "Name4"},
	}
	expected := `[{"ID":"1","Name":"Name1"},{"ID":"2","Name":"Name2"},{"ID":"3","Name":"Name3"},{"ID":"4","Name":"Name4"}]`

	// When
	s := csvx.ParserString[Struct](rows)

	// Then
	data, _ := json.Marshal(s)
	if string(data) != expected {
		t.Error("Parse csv format to array struct error", data)
	}
}

func TestParserStringEmpty(t *testing.T) {
	if result := csvx.ParserString[Struct](nil); len(result) != 0 {
		t.Error("ParserString empty error", result)
	}
}

func TestParser(t *testing.T) {
	// Given
	rows := [][]string{
		{"\ufeffID", "Name Space", "Age"},
		{"1", "Name1", "3.14"},
		{"2", "Name2", "3.14"},
		{"3", "Name3", "3.14"},
		{"4", "Name4", "3.14"},
	}
	jsonExpected := `[{"ID":1,"Name":"Name1","Age":3.14},{"ID":2,"Name":"Name2","Age":3.14},{"ID":3,"Name":"Name3","Age":3.14},{"ID":4,"Name":"Name4","Age":3.14}]`
	csvExpected := csvx.Utf8BOM + `"ID","Name Space","Age"
"1","Name1","3.14"
"2","Name2","3.14"
"3","Name3","3.14"
"4","Name4","3.14"`

	// When
	s := csvx.Parser[StructType](rows)
	c := csvx.Convert[StructType](s)

	// Then
	data, _ := json.Marshal(s)
	if string(data) != jsonExpected {
		t.Error("Parse csv format to array struct error", string(data))
	}
	if c != csvExpected {
		t.Error("Convert struct format to csv error\n", c)
	}
}

func TestParserEmpty(t *testing.T) {
	if result := csvx.Parser[StructType](nil); len(result) != 0 {
		t.Error("Parser empty error", result)
	}
}

func TestParserAllTypes(t *testing.T) {
	// Given
	rows := [][]string{
		{"ID", "Count", "Age", "Active"},
		{"1", "2", "3.14", "true"},
		{"bad", "bad", "bad", "bad"},
	}
	expected := `[{"ID":1,"Count":2,"Age":3.14,"Active":true},{"ID":0,"Count":0,"Age":0,"Active":false}]`

	// When
	result := csvx.Parser[StructAllTypes](rows)

	// Then
	data, _ := json.Marshal(result)
	if string(data) != expected {
		t.Error("Parser all types error", string(data))
	}
}

func TestParserPointer(t *testing.T) {
	// Given
	rows := [][]string{
		{"\ufeffID", "Name Space", "Age"},
		{"1", "Name1", "N/A"},
		{"2", "Name2", ""},
		{"3", "Name3", "3.14"},
		{"4", "Name4", "3.14"},
	}
	jsonExpected := `[{"ID":1,"Name":"Name1","Age":null},{"ID":2,"Name":"Name2","Age":null},{"ID":3,"Name":"Name3","Age":3.14},{"ID":4,"Name":"Name4","Age":3.14}]`
	csvExpected := csvx.Utf8BOM + `"ID","Name Space","Age"
"1","Name1",""
"2","Name2",""
"3","Name3","3.14"
"4","Name4","3.14"`

	// When
	s := csvx.Parser[StructPointerType](rows)
	c := csvx.Convert[StructPointerType](s)

	// Then
	data, _ := json.Marshal(s)
	if string(data) != jsonExpected {
		t.Error("Parse csv format to array struct error", string(data))
	}
	if c != csvExpected {
		t.Error("Convert struct format to csv error", c)
	}
}

func TestParserAllPointerTypes(t *testing.T) {
	// Given
	rows := [][]string{
		{"ID", "Count", "Age", "Active"},
		{"1", "2", "3.14", "true"},
		{"bad", "bad", "bad", "bad"},
	}
	expected := `[{"ID":1,"Count":2,"Age":3.14,"Active":true},{"ID":null,"Count":null,"Age":null,"Active":false}]`

	// When
	result := csvx.Parser[StructAllPointerTypes](rows)

	// Then
	data, _ := json.Marshal(result)
	if string(data) != expected {
		t.Error("Parser all pointer types error", string(data))
	}
}

func TestParserStructFieldPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Parser struct field should panic")
		}
	}()

	_ = csvx.Parser[StructWithStructField]([][]string{
		{"Value"},
		{"value"},
	})
}

func TestParserPointerStructFieldPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Parser pointer struct field should panic")
		}
	}()

	_ = csvx.Parser[StructWithPointerStructField]([][]string{
		{"Value"},
		{"value"},
	})
}

func TestParserFunc(t *testing.T) {
	// Given
	rows := [][]string{
		{"\ufeffID", "Name Space"},
		{"1", "Name1"},
		{"2", "Name2"},
		{"3", "Name3"},
		{"4", "Name4"},
	}

	// When
	called := 0
	_ = csvx.ParserFunc(true, rows, func(record []string) error {
		called++
		if record[0] == "ID" {
			t.Error("ParserFunc should exclude header")
		}
		return nil
	})

	// Then
	if called != 4 {
		t.Error("ParserFunc called error", called)
	}
}

func TestParserFuncIncludeHeaderAndError(t *testing.T) {
	// Given
	rows := [][]string{
		{"ID"},
		{"1"},
	}
	expectedErr := errors.New("stop")
	called := 0

	// When
	err := csvx.ParserFunc(false, rows, func(record []string) error {
		called++
		if called == 1 && record[0] != "ID" {
			t.Error("ParserFunc should include header")
		}
		return expectedErr
	})

	// Then
	if !errors.Is(err, expectedErr) {
		t.Error("ParserFunc error mismatch", err)
	}
	if called != 1 {
		t.Error("ParserFunc should stop on error")
	}
}

func TestParserByReader(t *testing.T) {
	// Given
	reader := csv.NewReader(strings.NewReader("ID;Name Space;Age\n1;Name1;3.14\n"))

	// When
	result := csvx.ParserByReader[StructType](reader, ';')

	// Then
	data, _ := json.Marshal(result)
	if string(data) != `[{"ID":1,"Name":"Name1","Age":3.14}]` {
		t.Error("ParserByReader error", string(data))
	}
}
