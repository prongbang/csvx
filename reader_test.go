package csvx_test

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/prongbang/csvx"
)

func TestReadByte(t *testing.T) {
	// Given
	path := t.TempDir() + "/test.csv"
	expected := []byte("ID,Name\n1,N1\n")
	if err := csvx.Append(path, []string{"ID", "Name"}); err != nil {
		t.Fatal(err)
	}
	if err := csvx.Append(path, []string{"1", "N1"}); err != nil {
		t.Fatal(err)
	}

	// When
	result := csvx.ReadByte(path)

	// Then
	if string(result) != string(expected) {
		t.Error("ReadByte error:", string(result))
	}
}

func TestReadByteMissingFile(t *testing.T) {
	if result := csvx.ReadByte(t.TempDir() + "/missing.csv"); len(result) != 0 {
		t.Error("ReadByte missing file error:", string(result))
	}
}

func TestByteReader(t *testing.T) {
	// Given
	data := []byte("#comment\nID;Name\n1;N1\n")

	// When
	result := csvx.ByteReader(data, func(r *csv.Reader) {
		r.Comma = ';'
	})

	// Then
	if len(result) != 2 || result[0][0] != "ID" || result[1][1] != "N1" {
		t.Error("ByteReader error:", result)
	}
}

func TestReader(t *testing.T) {
	// Given
	reader := csv.NewReader(strings.NewReader("#comment\nID,Name\n1,\"N1\n"))

	// When
	result := csvx.Reader(reader)

	// Then
	if len(result) != 2 || result[0][0] != "ID" || result[1][1] != "N1\n" {
		t.Error("Reader error:", result)
	}
}
