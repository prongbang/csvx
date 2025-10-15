package csvx_test

import (
	"encoding/json"
	"fmt"
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
	_ = csvx.ParserFunc(true, rows, func(record []string) error {
		fmt.Println(record)
		return nil
	})
}
