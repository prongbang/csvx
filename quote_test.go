package csvx_test

import (
	"bytes"
	"encoding/csv"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/prongbang/csvx"
)

func TestQuoteCell(t *testing.T) {
	cases := map[string]string{
		"":             `""`,
		"N1":           `"N1"`,
		"a,b":          `"a,b"`,
		`say "hi"`:     `"say ""hi"""`,
		"line1\nline2": "\"line1\nline2\"",
		"โลตัส, บางนา": `"โลตัส, บางนา"`,
	}
	for in, expected := range cases {
		if actual := csvx.QuoteCell(in); actual != expected {
			t.Errorf("QuoteCell(%q) = %q, expected %q", in, actual, expected)
		}
	}
}

func TestFormatQuoted(t *testing.T) {
	actual := csvx.FormatQuoted([]string{"1", "a,b", `x"y`, ""})
	expected := `"1","a,b","x""y",""`
	if actual != expected {
		t.Errorf("Expected %s, Actual %s", expected, actual)
	}
}

func TestQuotedWriter(t *testing.T) {
	// Given
	records := [][]string{
		{"ID", "Name"},
		{"1", "a,b"},
		{"2", `say "hi"`},
		{"3", "line1\nline2"},
	}
	var buf bytes.Buffer
	w := csvx.NewQuotedWriter(&buf)

	// When
	for _, r := range records {
		if err := w.Write(r); err != nil {
			t.Fatal(err)
		}
	}
	w.Flush()

	// Then
	if err := w.Error(); err != nil {
		t.Fatal(err)
	}
	expected := "\"ID\",\"Name\"\n\"1\",\"a,b\"\n\"2\",\"say \"\"hi\"\"\"\n\"3\",\"line1\nline2\"\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, Actual %q", expected, buf.String())
	}

	// A strict RFC 4180 reader gets the original records back.
	parsed, err := csv.NewReader(strings.NewReader(buf.String())).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed, records) {
		t.Errorf("Expected %v, Actual %v", records, parsed)
	}
}

func TestQuotedWriterCRLFAndWriteAll(t *testing.T) {
	var buf bytes.Buffer
	w := csvx.NewQuotedWriter(&buf)
	w.UseCRLF = true

	if err := w.WriteAll([][]string{{"a", "b"}, {"1", "2"}}); err != nil {
		t.Fatal(err)
	}

	expected := "\"a\",\"b\"\r\n\"1\",\"2\"\r\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, Actual %q", expected, buf.String())
	}
}

type failWriter struct{}

func (failWriter) Write([]byte) (int, error) { return 0, errors.New("disk full") }

func TestQuotedWriterError(t *testing.T) {
	w := csvx.NewQuotedWriter(failWriter{})
	_ = w.Write([]string{"a"})
	w.Flush()
	if w.Error() == nil {
		t.Error("Expected error from the underlying writer, got nil")
	}
}

type MyStructQuote struct {
	ID   int    `header:"ID" no:"1"`
	Name string `header:"Name" no:"2"`
}

func TestConvertEscapesDoubleQuote(t *testing.T) {
	m := []MyStructQuote{{ID: 1, Name: `say "hi", ok`}}
	expected := csvx.Utf8BOM + `"ID","Name"
"1","say ""hi"", ok"`

	if actual := csvx.Convert(m); actual != expected {
		t.Errorf("Expected %s, Actual %s", expected, actual)
	}
	if actual := csvx.ConvertWith(m); actual != expected {
		t.Errorf("ConvertWith: Expected %s, Actual %s", expected, actual)
	}
}

func TestTryConvertEscapesDoubleQuote(t *testing.T) {
	m := []MyStructQuote{{ID: 1, Name: `say "hi", ok`}}
	expected := csvx.Utf8BOM + `ID,Name
"1","say ""hi"", ok"
`

	if actual := csvx.TryConvert(m); actual != expected {
		t.Errorf("Expected %s, Actual %s", expected, actual)
	}
}
