package csvx_test

import (
	"fmt"
	"testing"

	"github.com/prongbang/csvx"
)

type MyStruct struct {
	Name  string `json:"name" db:"name" header:"Name Space" no:"2"`
	ID    int    `json:"id" db:"id" header:"ID" no:"1"`
	Other string
}

type MyStructPointer struct {
	Name    string  `json:"name" db:"name" header:"Name Space" no:"2"`
	ID      int     `json:"id" db:"id" header:"ID" no:"1"`
	Address *string `json:"address" db:"address" header:"Address" no:"4" default:"N/A"`
	Phone   *string `json:"phone" db:"phone" header:"Phone" no:"3" default:"NULL"`
	Email   *string `json:"email" db:"email" header:"Email" no:"5"`
	Age     *string `json:"age" db:"age" header:"Age" no:"6" default:"999"`
	Other   string
}

type MyStructDecimal struct {
	ID     int     `json:"id" header:"ID" no:"1"`
	Amount float64 `json:"amount" header:"Amount" no:"2" decimal:"1"`
	Price  float64 `json:"price" header:"Price" no:"3"`
	Rate   float32 `json:"rate" header:"Rate" no:"4"`
}

func TestConvert(t *testing.T) {
	// Given
	m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
	expected := csvx.Utf8BOM + `"ID","Name Space"
"1","N1"
"2","N2"`

	// When
	result := csvx.Convert[MyStruct](m)

	// Then
	if result != expected {
		t.Error("Convert error:", result)
	}
}

func TestTryConvertPointer(t *testing.T) {
	// Given
	age := "100"
	phone := "0876"
	m := []MyStructPointer{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2", Age: &age, Phone: &phone}}
	expected1 := csvx.Utf8BOM + `"ID","Name Space","Phone","Address","Email","Age"
"1","N1","NULL","N/A","","999"
"2","N2","0876","N/A","","100"`
	expected2 := csvx.Utf8BOM + `ID,Name Space,Phone,Address,Email,Age
"1","N1","NULL","N/A","","999"
"2","N2","0876","N/A","","100"
`

	// When
	result1 := csvx.Convert(m)
	result2 := csvx.TryConvert(m)

	// Then
	if result1 != expected1 {
		t.Error("Result1 error:\nexpected:", expected1, "\nactual:", result1)
	}
	if result2 != expected2 {
		t.Error("Result2 error:\nexpected:", expected2, "\nactual:", result2)
	}
}

func TestConvertWithDecimalTag(t *testing.T) {
	// Given
	m := []MyStructDecimal{{ID: 1, Amount: 1.23, Price: 2.345, Rate: 3.456}}
	expected := csvx.Utf8BOM + `"ID","Amount","Price","Rate"
"1","1.2","2.345","3.4560000896453857"`

	// When
	result := csvx.ConvertWith(m)

	// Then
	if result != expected {
		t.Error("ConvertWith decimal tag error:\nexpected:", expected, "\nactual:", result)
	}
}

func TestConvertWithDecimalOption(t *testing.T) {
	// Given
	m := []MyStructDecimal{{ID: 1, Amount: 1.23, Price: 2.345, Rate: 3.456}}
	expected := csvx.Utf8BOM + `"ID","Amount","Price","Rate"
"1","1","2.35","3.5"`

	// When
	result := csvx.ConvertWith(
		m,
		csvx.WithDecimal("Amount", 0),
		csvx.WithDecimal("Price", 2),
		csvx.WithDecimals(map[string]int{
			"Rate": 1,
		}),
	)

	// Then
	if result != expected {
		t.Error("ConvertWith decimal option error:\nexpected:", expected, "\nactual:", result)
	}
}

func TestConvertWithDecimalsOnly(t *testing.T) {
	// Given
	m := []MyStructDecimal{{ID: 1, Amount: 1.23, Price: 2.345, Rate: 3.456}}
	expected := csvx.Utf8BOM + `"ID","Amount","Price","Rate"
"1","1.2","2.345","3.4560000896453857"`

	// When
	result := csvx.ConvertWith(
		m,
		csvx.WithDecimals(map[string]int{
			"Amount": 1,
		}),
	)

	// Then
	if result != expected {
		t.Error("ConvertWith decimals only error:\nexpected:", expected, "\nactual:", result)
	}
}

func TestConvertWithDecimalOptionByHeaderAndInvalidValues(t *testing.T) {
	// Given
	type DecimalByHeader struct {
		Amount float64 `json:"amount" header:"Total Amount" no:"1" decimal:"bad"`
		Price  float64 `json:"price" header:"Price" no:"2" decimal:"-1"`
	}
	m := []DecimalByHeader{{Amount: 1.23, Price: 2.345}}
	expected := csvx.Utf8BOM + `"Total Amount","Price"
"1.2","2.345"`

	// When
	result := csvx.ConvertWith(
		m,
		csvx.WithDecimal("Total Amount", 1),
		csvx.WithDecimal("Ignored", -1),
		csvx.WithDecimals(map[string]int{
			"Price":   -2,
			"Missing": 3,
		}),
	)

	// Then
	if result != expected {
		t.Error("ConvertWith decimal header option error:\nexpected:", expected, "\nactual:", result)
	}
}

func TestConvertWithIgnoreDoubleQuote(t *testing.T) {
	// Given
	m := []MyStruct{{ID: 1, Name: "N1"}}
	expected := csvx.Utf8BOM + `ID,Name Space
1,N1`

	// When
	result := csvx.ConvertWith(m, csvx.WithIgnoreDoubleQuote(true))

	// Then
	if result != expected {
		t.Error("ConvertWith ignore double quote error:\nexpected:", expected, "\nactual:", result)
	}
}

func TestConvertEmpty(t *testing.T) {
	if result := csvx.Convert[MyStruct](nil); result != "" {
		t.Error("Convert empty error:", result)
	}
	if result := csvx.ConvertWith[MyStruct](nil); result != "" {
		t.Error("ConvertWith empty error:", result)
	}
}

func TestConvertIgnoreDoubleQuote(t *testing.T) {
	// Given
	m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
	expected := csvx.Utf8BOM + `ID,Name Space
1,N1
2,N2`

	// When
	result := csvx.Convert[MyStruct](m, true)

	// Then
	if result != expected {
		t.Error("Convert error:", result)
	}
}

func TestManualConvert(t *testing.T) {
	// Given
	m := []MyStruct{{ID: 1, Name: "N1"}}
	expected := csvx.Utf8BOM + `ID,Name Space
1,N1
`

	// When
	result := csvx.ManualConvert[MyStruct](m,
		[]string{"ID", "Name Space"},
		func(data MyStruct) []string {
			return []string{
				fmt.Sprintf("%d", data.ID),
				data.Name,
			}
		},
	)

	// Then
	if result != expected {
		t.Error("ManualConvert error:\nexpected:", expected, "\nactual:", result)
	}
}

func TestManualConvertEmpty(t *testing.T) {
	if result := csvx.ManualConvert[MyStruct](nil, []string{"ID"}, func(data MyStruct) []string {
		return []string{fmt.Sprintf("%d", data.ID)}
	}); result != "" {
		t.Error("ManualConvert empty error:", result)
	}
}

func TestTryConvertIgnoreDoubleQuoteAndInvalidNo(t *testing.T) {
	// Given
	type InvalidNoStruct struct {
		ID   int    `header:"ID" no:"1"`
		Name string `header:"Name" no:"bad"`
	}
	m := []InvalidNoStruct{{ID: 1, Name: "N1"}}
	expected := csvx.Utf8BOM + `ID
1
`

	// When
	result := csvx.TryConvert(m, true)

	// Then
	if result != expected {
		t.Error("TryConvert invalid no error:\nexpected:", expected, "\nactual:", result)
	}
}

func TestTryConvertFloat(t *testing.T) {
	// Given
	m := []MyStructDecimal{{ID: 1, Amount: 1.23, Price: 2.345, Rate: 3.456}}
	expected := csvx.Utf8BOM + `ID,Amount,Price,Rate
"1","1.23","2.345","3.4560000896453857"
`

	// When
	result := csvx.TryConvert(m)

	// Then
	if result != expected {
		t.Error("TryConvert float error:\nexpected:", expected, "\nactual:", result)
	}
}

func TestTryConvertEmpty(t *testing.T) {
	if result := csvx.TryConvert[MyStruct](nil); result != "" {
		t.Error("TryConvert empty error:", result)
	}
}

func BenchmarkConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
		expected := csvx.Utf8BOM + `"ID","Name Space"
"1","N1"
"2","N2"`

		// When
		result := csvx.Convert[MyStruct](m)

		// Then
		if result != expected {
			b.Error("Convert error:", result)
		}
	}
}

func BenchmarkConvertWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStructDecimal{{ID: 1, Amount: 1.23, Price: 2.345, Rate: 3.456}}
		expected := csvx.Utf8BOM + `"ID","Amount","Price","Rate"
"1","1","2.35","3.5"`

		// When
		result := csvx.ConvertWith(
			m,
			csvx.WithDecimal("Amount", 0),
			csvx.WithDecimal("Price", 2),
			csvx.WithDecimals(map[string]int{
				"Rate": 1,
			}),
		)

		// Then
		if result != expected {
			b.Error("ConvertWith error:", result)
		}
	}
}

func BenchmarkManualConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
		expected := csvx.Utf8BOM + `ID,Name Space
1,N1
2,N2
`

		// When
		result := csvx.ManualConvert[MyStruct](m,
			[]string{"ID", "Name Space"},
			func(data MyStruct) []string {
				return []string{
					fmt.Sprintf("%d", data.ID),
					data.Name,
				}
			},
		)

		// Then
		if result != expected {
			b.Error("Convert error:", result)
		}
	}
}

func BenchmarkTryConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Given
		m := []MyStruct{{ID: 1, Name: "N1"}, {ID: 2, Name: "N2"}}
		expected := csvx.Utf8BOM + `ID,Name Space
"1","N1"
"2","N2"
`

		// When
		result := csvx.TryConvert(m)

		// Then
		if result != expected {
			b.Error("Convert error:", result)
		}
	}
}
