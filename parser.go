package csvx

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"strconv"
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

// ParserString is a generic function that takes a slice of slices of strings as input and returns a slice of values of type T,
// where T is a type parameter that represents the desired output type. The input slice should represent a CSV file
// or other tabular data in which each inner slice represents a single row of data, and each element in the inner slice represents
// a single field value. This function will attempt to parse each field value into the corresponding type T using the built-in strconv package.
// If parsing fails or the input slice is empty, an empty slice of type T will be returned.
//
//	type Struct struct {
//		  ID   string `header:"ID"`
//		  Name string `header:"Name Space"`
//	}
//
//	rows := [][]string{
//	   {"ID", "Name Space"},
//	   {"1", "Name1"},
//	}
//
// s := csvx.ParserString[Struct](rows)
func ParserString[T any](rows [][]string) []T {
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
				fieldTag := f.Tag.Get("header")
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

// Parser parses the provided input data and returns the result.
// It handles different formats based on the input type.
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
				fieldTag := f.Tag.Get("header")
				head := RemoveDoubleQuote(header[j])
				return fieldTag == fmt.Sprintf("%v", head)
			})

			if structField.IsValid() {
				// Convert the value based on the field kind
				switch structField.Kind() {
				case reflect.Ptr:
					// Handle pointer types
					fieldType := structField.Type()
					elemType := fieldType.Elem()
					ptrValue := reflect.New(elemType)
					switch elemType.Kind() {
					case reflect.String:
						ptrValue.Elem().SetString(field)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						value, err := strconv.ParseInt(field, 10, 64)
						if err == nil {
							ptrValue.Elem().SetInt(value)
						} else {
							structField.Set(reflect.Zero(fieldType))
							continue
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						value, err := strconv.ParseUint(field, 10, 64)
						if err == nil {
							ptrValue.Elem().SetUint(value)
						} else {
							structField.Set(reflect.Zero(fieldType))
							continue
						}
					case reflect.Float32, reflect.Float64:
						value, err := strconv.ParseFloat(field, 64)
						if err == nil {
							ptrValue.Elem().SetFloat(value)
						} else {
							structField.Set(reflect.Zero(fieldType))
							continue
						}
					case reflect.Bool:
						value, err := strconv.ParseBool(field)
						if err == nil {
							ptrValue.Elem().SetBool(value)
						}
					case reflect.Struct:
						ptrValue.Elem().Set(reflect.ValueOf(field))
					}
					structField.Set(ptrValue)
				default:
					// Handle non-pointer types as before
					switch structField.Kind() {
					case reflect.String:
						structField.SetString(field)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						value, err := strconv.ParseInt(field, 10, 64)
						if err == nil {
							structField.SetInt(value)
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						value, err := strconv.ParseUint(field, 10, 64)
						if err == nil {
							structField.SetUint(value)
						}
					case reflect.Float32, reflect.Float64:
						value, err := strconv.ParseFloat(field, 64)
						if err == nil {
							structField.SetFloat(value)
						}
					case reflect.Bool:
						value, err := strconv.ParseBool(field)
						if err == nil {
							structField.SetBool(value)
						}
					case reflect.Struct:
						structField.Set(reflect.ValueOf(field))
					}
				}
			}
		}

		structs = append(structs, record.Data)
	}

	return structs
}

// ParserFunc processes the input data using a custom parsing function.
// This allows for flexible and reusable parsing logic.
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

// ParserByReader parses data from an io.Reader and returns the result.
// This is useful for streaming data or reading from large files.
func ParserByReader[T any](ir *csv.Reader, delimiter ...rune) []T {
	d := ','
	if len(delimiter) > 0 {
		d = delimiter[0]
	}
	return Parser[T](Reader(ir, func(r *csv.Reader) {
		r.Comma = d
	}))
}
