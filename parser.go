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

// FileHeaderReader
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

// Parser
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
				return fieldTag == fmt.Sprintf("%v", header[j])
			})

			if structField.IsValid() {
				structField.SetString(field)
			}
		}

		structs = append(structs, record.Data)
	}

	return structs
}
