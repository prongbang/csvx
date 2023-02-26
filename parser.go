package csvx

import (
	"fmt"
	"reflect"
)

type model[T any] struct {
	Data T
}

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
