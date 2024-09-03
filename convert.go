package csvx

import (
	"encoding/csv"
	"fmt"
	"reflect"
	"sort"
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
	if size == 0 {
		return ""
	}

	// Config format value
	valueFormatCore := "%v"
	valueFormat := "\"%v\""
	if len(ignoreDoubleQuote) > 0 {
		valueFormat = valueFormatCore
	}

	// Initialize the element
	var headers strings.Builder
	rows := make([]strings.Builder, size)

	first := data[0]
	fel := reflect.ValueOf(&first).Elem()
	numField := fel.NumField()
	for c := 0; c < numField; c++ {
		field, fOk := headerLookup[T](first, c)
		index, iOk := noLookup[T](first, c)
		if fOk || iOk {
			fmt.Println(field, index)
		}
	}

	// Mapping
	sheets := []string{}
	for r, d := range data {
		el := reflect.ValueOf(&d).Elem()

		colsRaw := el.NumField()
		for c := 0; c < colsRaw; c++ {
			value := el.Field(c)
			field, fOk := headerLookup[T](d, c)
			no, iOk := noLookup[T](d, c)
			if !fOk || !iOk {
				continue
			}

			if _, err := strconv.Atoi(no); err == nil {
				isNotLast := c < colsRaw-1
				if r == 0 {
					headers.WriteString(fmt.Sprintf(valueFormat, field))
					if isNotLast {
						headers.WriteString(",")
					}
				}
				if IsFloat(value.Type()) {
					rows[r].WriteString(fmt.Sprintf(valueFormat, F64ToString(value.Float())))
				} else {
					nValue := ""
					if IsPointer(value.Type()) {
						if value.Elem().IsValid() {
							nValue = RemoveDoubleQuote(fmt.Sprintf(valueFormatCore, value.Elem()))
						}
					} else {
						nValue = RemoveDoubleQuote(fmt.Sprintf(valueFormatCore, value))
					}
					rows[r].WriteString(fmt.Sprintf(valueFormat, nValue))
				}

				if isNotLast {
					rows[r].WriteString(",")
				}
			}
		}

		// Convert array to csv format
		if len(sheets) == 0 {
			sheets = append(sheets, headers.String())
		}
		sheets = append(sheets, rows[r].String())
	}

	// Add enter end line
	result := strings.Join(sheets, "\n")
	return result
}

func ConvertToCSV[T any](data []T) string {
	if len(data) == 0 {
		return ""
	}

	// Use reflection to get the type of the struct
	t := reflect.TypeOf(data[0])

	// Create a slice to hold the fields and their `no` tags
	fields := make([]struct {
		field reflect.StructField
		no    int
	}, t.NumField())

	// Populate the fields slice
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		noTag := field.Tag.Get("no")
		no := 0
		fmt.Sscanf(noTag, "%d", &no)
		fields[i] = struct {
			field reflect.StructField
			no    int
		}{field, no}
	}

	// Sort fields by `no` tag
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].no < fields[j].no
	})

	// Create a CSV writer
	var csvBuilder strings.Builder
	csvWriter := csv.NewWriter(&csvBuilder)

	// Write headers
	headers := make([]string, len(fields))
	for i, f := range fields {
		headers[i] = f.field.Tag.Get("header")
	}
	csvWriter.Write(headers)

	// Write data
	for _, record := range data {
		values := make([]string, len(fields))
		val := reflect.ValueOf(record)
		for i, f := range fields {
			values[i] = fmt.Sprintf("%v", val.FieldByName(f.field.Name))
		}
		csvWriter.Write(values)
	}

	csvWriter.Flush()
	return csvBuilder.String()
}

func ConvertToCSV2[T any](data []T) string {
	if len(data) == 0 {
		return ""
	}

	// Use reflection to get the type of the struct
	t := reflect.TypeOf(data[0])

	// Initialize arrays to hold headers and field names in the correct order
	headers := make([]string, t.NumField())
	fieldNames := make([]string, t.NumField())

	// Populate the arrays according to the `no` tag
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		noTag := field.Tag.Get("no")
		no := 0
		fmt.Sscanf(noTag, "%d", &no)
		// Place headers and fields in the correct position based on `no` tag
		headers[no-1] = field.Tag.Get("header")
		fieldNames[no-1] = field.Name
	}

	// Create a CSV writer
	var csvBuilder strings.Builder
	csvWriter := csv.NewWriter(&csvBuilder)

	// Write headers
	csvWriter.Write(headers)

	// Write data
	for _, record := range data {
		values := make([]string, len(fieldNames))
		val := reflect.ValueOf(record)
		for i, fieldName := range fieldNames {
			values[i] = fmt.Sprintf("%v", val.FieldByName(fieldName))
		}
		csvWriter.Write(values)
	}

	csvWriter.Flush()
	return csvBuilder.String()
}

func TryConvert[T any](data []T, ignoreDoubleQuote ...bool) string {
	size := len(data)
	if size == 0 {
		return ""
	}

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

func headerLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("header")
}

func noLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("no")
}
