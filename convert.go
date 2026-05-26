package csvx

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Utf8BOM represents the Byte Order Mark (BOM) for UTF-8 encoding.
// It is used to signal UTF-8 encoding in text files for better compatibility with tools like Excel.
const Utf8BOM = "\uFEFF"

type convertConfig struct {
	ignoreDoubleQuote bool
	decimals          map[string]int
}

type ConvertOption func(*convertConfig)

// WithIgnoreDoubleQuote controls whether ConvertWith wraps CSV cells with double quotes.
func WithIgnoreDoubleQuote(ignore bool) ConvertOption {
	return func(config *convertConfig) {
		config.ignoreDoubleQuote = ignore
	}
}

// WithDecimal sets decimal precision for a float field by struct field name, field tag, or header tag.
func WithDecimal(field string, precision int) ConvertOption {
	return func(config *convertConfig) {
		if precision < 0 {
			return
		}
		if config.decimals == nil {
			config.decimals = map[string]int{}
		}
		config.decimals[field] = precision
	}
}

// WithDecimals sets decimal precision for multiple float fields.
func WithDecimals(decimals map[string]int) ConvertOption {
	return func(config *convertConfig) {
		for field, precision := range decimals {
			if precision < 0 {
				continue
			}
			if config.decimals == nil {
				config.decimals = map[string]int{}
			}
			config.decimals[field] = precision
		}
	}
}

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
	config := convertConfig{}
	if len(ignoreDoubleQuote) > 0 {
		config.ignoreDoubleQuote = ignoreDoubleQuote[0]
	}
	return convert(data, config)
}

// ConvertWith converts array struct to csv format with options.
func ConvertWith[T any](data []T, options ...ConvertOption) string {
	config := convertConfig{}
	for _, option := range options {
		option(&config)
	}
	return convert(data, config)
}

func convert[T any](data []T, config convertConfig) string {
	size := len(data)
	if size > 0 {

		// Config format value
		valueFormatCore := "%v"
		valueFormat := "\"%v\""
		if config.ignoreDoubleQuote {
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
						rows[r][i-1] = fmt.Sprintf(valueFormat, floatToString(d, c, value.Float(), config))
					} else {
						nValue := ""
						if IsPointer(value.Type()) {
							if value.Elem().IsValid() {
								elem := value.Elem()
								if IsFloat(elem.Type()) {
									nValue = floatToString(d, c, elem.Float(), config)
								} else {
									nValue = RemoveDoubleQuote(fmt.Sprintf(valueFormatCore, elem))
								}
							} else {
								// Use default from tag
								if def, dOk := defaultLookup[T](d, c); dOk {
									nValue = def
								}
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
		return Utf8BOM + result
	}

	return ""
}

func decimalLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("decimal")
}

func floatToString[T any](d T, c int, value float64, config convertConfig) string {
	precision, ok := decimalPrecision(d, c, config)
	if !ok {
		return F64ToString(value)
	}
	return strconv.FormatFloat(value, 'f', precision, 64)
}

func decimalPrecision[T any](d T, c int, config convertConfig) (int, bool) {
	structField := reflect.ValueOf(d).Type().Field(c)
	if config.decimals != nil {
		if precision, ok := config.decimals[structField.Name]; ok {
			return precision, true
		}
		for _, tag := range []string{"field", "header"} {
			if field, ok := structField.Tag.Lookup(tag); ok {
				if precision, ok := config.decimals[field]; ok {
					return precision, true
				}
			}
		}
	}

	decimal, ok := decimalLookup(d, c)
	if !ok {
		return 0, false
	}

	precision, err := strconv.Atoi(decimal)
	if err != nil || precision < 0 {
		return 0, false
	}
	return precision, true
}

// ManualConvert performs a manual conversion of the input data.
// It applies the specified rules or transformations to achieve the desired output.
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
	return fmt.Sprintf("%s%s", Utf8BOM, buffer.String())
}

// TryConvert attempts to convert the input data to the specified format.
// It handles errors gracefully and returns the converted result along with an error (if any).
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
	dmap := make(map[int]string)
	for i := 0; i < numField; i++ {
		header, hOk := headerLookup[T](data[0], i)
		noStr, nOk := noLookup[T](data[0], i)
		def, _ := defaultLookup[T](data[0], i)
		if hOk && nOk {
			no, err := strconv.Atoi(noStr)
			if err != nil {
				continue
			}
			// Convert to index of array
			index := no - 1
			nmap[index] = i
			hmap[index] = header
			dmap[index] = def
			cols++
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
					} else {
						// Use default from tag
						value = dmap[c]
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

	return fmt.Sprintf("%s%s\n%s", Utf8BOM, headers.String(), records.String())
}

func headerLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("header")
}

func noLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("no")
}

func defaultLookup[T any](d T, c int) (string, bool) {
	return reflect.ValueOf(d).Type().Field(c).Tag.Lookup("default")
}
