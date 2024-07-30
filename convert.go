package csvx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Format formats a slice of strings as a single, comma-separated string. Each element in the slice will be separated
// by a comma and a space. This function can be used to generate formatted output for CSV files or other data formats
// that use comma-separated values. If the input slice is empty, this function will return an empty string.
func Format(cell []string) string {
	return strings.Join(cell, ",")
}

// Convert array struct to csv format
// Struct supported
//
//	type MyStruct struct {
//		Name string `json:"name" header:"Name" no:"2"`
//		ID   int    `json:"id" header:"ID" no:"1"`
//	}
//
// m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
// csv := csvx.Convert[MyStruct](m)
//
// Result:
//
//	"ID","Name"
//	"1","N1"
//	"2","N2"
func Convert[T any](data []T, ignoreDoubleQuote ...bool) string {
	size := len(data)
	if size > 0 {

		// Config format value
		valueFormatCore := "%v"
		valueFormat := "\"%v\""
		if len(ignoreDoubleQuote) > 0 {
			valueFormat = valueFormatCore
		}

		// Initialize the element
		var headers []string
		rows := make([][]string, size)

		// Mapping
		sheets := []string{}
		for r, d := range data {
			el := reflect.ValueOf(&d).Elem()

			colsRaw := el.NumField()
			cols := 0
			for c := 0; c < colsRaw; c++ {
				_, fOk := fieldLookup[T](d, c)
				_, iOk := indexLookup[T](d, c)
				if !fOk || !iOk {
					continue
				}
				cols++
			}
			if headers == nil {
				headers = make([]string, cols)
			}
			if len(rows[r]) == 0 {
				rows[r] = make([]string, cols)
			}

			for c := 0; c < colsRaw; c++ {
				value := el.Field(c)
				field, fOk := fieldLookup[T](d, c)
				index, iOk := indexLookup[T](d, c)
				if !fOk || !iOk {
					continue
				}

				if i, err := strconv.Atoi(index); err == nil {
					if r == 0 {
						headers[i-1] = fmt.Sprintf(valueFormat, field)
					}
					if IsFloat(value.Type()) {
						rows[r][i-1] = fmt.Sprintf(valueFormat, F64ToString(value.Float()))
					} else {
						nValue := ""
						if IsPointer(value.Type()) {
							if value.Elem().IsValid() {
								nValue = RemoveDoubleQuote(fmt.Sprintf(valueFormatCore, value.Elem()))
							}
						} else {
							nValue = RemoveDoubleQuote(fmt.Sprintf(valueFormatCore, value))
						}
						rows[r][i-1] = fmt.Sprintf(valueFormat, nValue)
					}
				}
			}

			// Convert array to csv format
			if len(sheets) == 0 {
				sheets = append(sheets, Format(headers))
			}
			sheets = append(sheets, Format(rows[r]))
		}

		// Add enter end line
		result := strings.Join(sheets, "\n")
		return result
	}

	return ""
}

func fieldLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("header")
}

func indexLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("no")
}
