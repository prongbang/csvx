package csvx

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
)

type model[T any] struct {
	Data T
}

// FileHeaderReader extracts the header of a multipart file specified by the given *multipart.FileHeader parameter
// and returns a slice of slices of strings representing the parsed header. Each slice in the result represents
// a single header field, where the first element is the header field name and the second element is the header field value.
// If the header is empty or cannot be parsed, an empty slice will be returned. If an error occurs during the operation,
// an error value will be returned.
// Ex:
// file, _ := c.FormFile("file")
// rows, err := csvx.FileHeaderReader(file)
func FileHeaderReader(fileHeader *multipart.FileHeader) ([][]string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return [][]string{}, err
	}

	// Parse the file
	r := csv.NewReader(bufio.NewReader(file))

	// Read the records
	_, err = r.Read()
	if err != nil {
		return [][]string{}, err
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

	return rows, nil
}

// Parser is a generic function that takes a slice of slices of strings as input and returns a slice of values of type T,
// where T is a type parameter that represents the desired output type. The input slice should represent a CSV file
// or other tabular data in which each inner slice represents a single row of data, and each element in the inner slice represents
// a single field value. This function will attempt to parse each field value into the corresponding type T using the built-in strconv package.
// If parsing fails or the input slice is empty, an empty slice of type T will be returned.
//
//	type Struct struct {
//		  ID   string `field:"ID"`
//		  Name string `field:"Name Space"`
//	}
//
//	rows := [][]string{
//	   {"ID", "Name Space"},
//	   {"1", "Name1"},
//	}
//
// s := csvx.Parser[Struct](rows)
func Parser[T any](rows [][]string) []T {
	var structs []T

	if len(rows) == 0 {
		return structs
	}

	header := rows[0]
	for i, row := range rows {
		if i == 0 {
			continue
		}

		record := model[T]{}
		structValue := reflect.ValueOf(&record.Data).Elem()

		for j, field := range row {
			structField := structValue.FieldByNameFunc(func(fieldName string) bool {
				f, _ := reflect.TypeOf(record.Data).FieldByName(fieldName)
				fieldTag := f.Tag.Get("field")
				head := RemoveDoubleQuote(header[j])
				return fieldTag == fmt.Sprintf("%v", head)
			})

			if structField.IsValid() {
				structField.SetString(field)
			}
		}

		structs = append(structs, record.Data)
	}

	return structs
}

// ParserFunc
//
//	err := csvx.ParserFunc(true, rows, func (record []string) {
//		return nil
//	})
func ParserFunc(excludeHeader bool, rows [][]string, onRecord func([]string) error) error {
	for i, row := range rows {
		if excludeHeader && i == 0 {
			continue
		}
		if err := onRecord(row); err != nil {
			return err
		}
	}
	return nil
}

func ParserByReader[T any](ir *csv.Reader, delimiter ...rune) []T {
	d := ','
	if len(delimiter) > 0 {
		d = delimiter[0]
	}
	return Parser[T](Reader(ir, func(r *csv.Reader) {
		r.Comma = d
	}))
}
