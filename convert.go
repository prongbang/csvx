package csvx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Format(cell []string) string {
	return strings.Join(cell, ",")
}

// Convert array struct to csv format
// Struct supported
//
//	type MyStruct struct {
//		Name string `json:"name" field:"Name" index:"2"`
//		ID   int    `json:"id" field:"ID" index:"1"`
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
		valueFormat := "\"%v\""
		if len(ignoreDoubleQuote) > 0 {
			valueFormat = "%v"
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
						rows[r][i-1] = fmt.Sprintf(valueFormat, value)
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
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("field")
}

func indexLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("index")
}
