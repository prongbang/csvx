package csvx_test

import (
	"os"
	"testing"

	"github.com/prongbang/csvx"
)

func TestAppendCreateAndAppend(t *testing.T) {
	// Given
	path := t.TempDir() + "/test.csv"

	// When
	err := csvx.Append(path, []string{"ID", "Name"})
	if err != nil {
		t.Fatal(err)
	}
	err = csvx.Append(path, []string{"1", "N1"})
	if err != nil {
		t.Fatal(err)
	}

	// Then
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	expected := "ID,Name\n1,N1\n"
	if string(data) != expected {
		t.Error("Append error:\nexpected:", expected, "\nactual:", string(data))
	}
}

func TestAppendInvalidPath(t *testing.T) {
	err := csvx.Append(t.TempDir()+"/missing/test.csv", []string{"ID"})
	if err == nil {
		t.Error("Append should return error")
	}
}

func TestAppendDirectoryPath(t *testing.T) {
	err := csvx.Append(t.TempDir(), []string{"ID"})
	if err == nil {
		t.Error("Append directory path should return error")
	}
}
