package csvx

import (
	"bytes"
	"encoding/csv"
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
				_, fOk := headerLookup[T](d, c)
				_, iOk := noLookup[T](d, c)
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
				field, fOk := headerLookup[T](d, c)
				index, iOk := noLookup[T](d, c)
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

func ManualConvert[T any](data []T, headers []string, onRecord func(data T) []string) string {
	size := len(data)
	if size == 0 {
		return ""
	}

	var buffer bytes.Buffer
	w := csv.NewWriter(&buffer)

	_ = w.Write(headers)
	for _, d := range data {
		row := onRecord(d)
		_ = w.Write(row)
	}
	w.Flush()
	return buffer.String()
}

func TryConvert[T any](data []T, ignoreDoubleQuote ...bool) string {
	if len(data) == 0 {
		return ""
	}

	// Config format value
	valueFormatCore := "%v"
	valueFormat := "\"%v\""
	if len(ignoreDoubleQuote) > 0 {
		valueFormat = valueFormatCore
	}

	// Use reflection to get the type of the struct
	t := reflect.TypeOf(data[0])

	cols := 0
	numField := t.NumField()
	hmap := make(map[int]string)
	nmap := make(map[int]int)
	for i := 0; i < numField; i++ {
		header, hOk := headerLookup[T](data[0], i)
		noStr, nOk := noLookup[T](data[0], i)
		if hOk && nOk {
			no, err := strconv.Atoi(noStr)
			if err != nil {
				continue
			}
			// Convert to index of array
			index := no - 1
			nmap[index] = i
			hmap[index] = header
			cols += 1
		}
	}

	var headers strings.Builder
	var records strings.Builder

	for r, d := range data {
		el := reflect.ValueOf(&d).Elem()
		for c := 0; c < cols; c++ {
			idx := nmap[c]

			// Header
			header := hmap[c]
			if r == 0 {
				headers.WriteString(fmt.Sprintf("%v", header))
				if c < cols-1 {
					headers.WriteString(",")
				}
			}

			// Records
			field := el.Field(idx)

			if IsFloat(field.Type()) {
				records.WriteString(fmt.Sprintf(valueFormat, F64ToString(field.Float())))
			} else {
				value := ""
				if IsPointer(field.Type()) {
					if field.Elem().IsValid() {
						value = fmt.Sprintf(valueFormatCore, field.Elem())
					}
				} else {
					value = fmt.Sprintf(valueFormatCore, field)
				}
				records.WriteString(fmt.Sprintf(valueFormat, value))
			}
			if c < cols-1 {
				records.WriteString(",")
			} else {
				records.WriteString("\n")
			}
		}
	}

	return fmt.Sprintf("%s\n%s", headers.String(), records.String())
}

func headerLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("header")
}

func noLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("no")
}
