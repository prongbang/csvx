package csvx

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
)

// ReadByte reads the entire contents of the file with the specified filename and returns them as a slice of bytes.
// If an error occurs during the operation, a nil slice will be returned. This function can be used to read the contents
// of text files or other files that are encoded as byte streams. Note that this function may not be suitable for reading
// large files, as it reads the entire file into memory at once. For large files, consider using the os package or
// a buffered reader to read the file in smaller chunks.
func ReadByte(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}
	}
	return data
}

// ByteReader creates an io.Reader from a byte slice.
// It allows the byte data to be read sequentially as a stream.
func ByteReader(data []byte, options ...func(r *csv.Reader)) [][]string {
	// Create a bytes.Reader from the byte slice
	byteReader := bytes.NewReader(data)

	// Parse the file
	r := csv.NewReader(byteReader)

	return Reader(r, options...)
}

// Reader wraps an existing io.Reader to provide additional functionality.
// It may include features like buffering or line-by-line reading.
func Reader(r *csv.Reader, options ...func(r *csv.Reader)) [][]string {
	r.LazyQuotes = true
	r.Comma = ','
	r.Comment = '#' // Set the comment character (lines beginning with this are ignored)

	for _, option := range options {
		option(r)
	}

	// Iterate through the records
	rows := [][]string{}
	for {
		// Read each record from csv
		record, e := r.Read()
		if e == io.EOF {
			break
		}
		rows = append(rows, record)
	}

	return rows
}
