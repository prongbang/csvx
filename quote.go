package csvx

import (
	"bufio"
	"io"
	"strings"
)

// QuoteCell wraps a cell with double quotes and escapes every double quote inside it as "" (RFC 4180),
// so a value containing a comma, a double quote or a line break stays in one column.
//
// QuoteCell(`say "hi", ok`) => `"say ""hi"", ok"`
func QuoteCell(cell string) string {
	return `"` + strings.ReplaceAll(cell, `"`, `""`) + `"`
}

// FormatQuoted formats a row as one CSV line with every cell wrapped by QuoteCell.
//
// FormatQuoted([]string{"1", "a,b"}) => `"1","a,b"`
func FormatQuoted(cells []string) string {
	quoted := make([]string, len(cells))
	for i, cell := range cells {
		quoted[i] = QuoteCell(cell)
	}
	return Format(quoted)
}

// QuotedWriter writes CSV records with every cell wrapped in double quotes, row by row.
// It mirrors encoding/csv.Writer (Write, WriteAll, Flush, Error), so it can replace csv.NewWriter
// where every cell must be quoted, including streaming into a gzip or file writer.
// encoding/csv.Writer only quotes the cells that need it.
type QuotedWriter struct {
	// UseCRLF ends each record with \r\n instead of \n.
	UseCRLF bool
	w       *bufio.Writer
}

// NewQuotedWriter returns a QuotedWriter that writes to w.
func NewQuotedWriter(w io.Writer) *QuotedWriter {
	return &QuotedWriter{w: bufio.NewWriter(w)}
}

// Write writes a single record. Writes are buffered, so Flush must be called to ensure the record reaches w.
func (q *QuotedWriter) Write(record []string) error {
	for i, cell := range record {
		if i > 0 {
			if err := q.w.WriteByte(','); err != nil {
				return err
			}
		}
		if _, err := q.w.WriteString(QuoteCell(cell)); err != nil {
			return err
		}
	}
	lineEnd := "\n"
	if q.UseCRLF {
		lineEnd = "\r\n"
	}
	_, err := q.w.WriteString(lineEnd)
	return err
}

// WriteAll writes multiple records followed by a Flush, returning any error from the Flush.
func (q *QuotedWriter) WriteAll(records [][]string) error {
	for _, record := range records {
		if err := q.Write(record); err != nil {
			return err
		}
	}
	return q.w.Flush()
}

// Flush writes any buffered data to the underlying io.Writer. Call Error to check whether it failed.
func (q *QuotedWriter) Flush() {
	_ = q.w.Flush()
}

// Error reports any error that has occurred during a previous Write or Flush.
func (q *QuotedWriter) Error() error {
	_, err := q.w.Write(nil)
	return err
}
