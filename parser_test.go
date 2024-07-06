package csvx_test

import (
	"encoding/json"
	"fmt"
	"github.com/prongbang/csvx"
	"testing"
)

type Struct struct {
	ID   string `field:"ID"`
	Name string `field:"Name Space"`
}

func TestParser(t *testing.T) {
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
	s := csvx.Parser[Struct](rows)

	// Then
	data, _ := json.Marshal(s)
	if string(data) != expected {
		t.Error("Parse csv format to array struct error", data)
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
