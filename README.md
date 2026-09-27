# CSVX

[![Go Coverage](https://github.com/prongbang/csvx/wiki/coverage.svg)](https://raw.githack.com/wiki/prongbang/csvx/coverage.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/prongbang/csvx)](https://goreportcard.com/report/github.com/prongbang/csvx)

Convert array struct to csv format and Parse csv format to array struct with Golang

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/prongbang)

## Install

```shell
go get github.com/prongbang/csvx
```

## Define struct for Convert

Add `header` for mapping in csv header and `no` start with 1 for sort header.
Use `default` for nil pointer values and `decimal` for fixed float precision.

```go
type MyStruct struct {
    Name    string  `json:"name" header:"Name Space" no:"2"`
    ID      int     `json:"id" header:"ID" no:"1"`
    Address *string `json:"address" header:"Address" no:"4" default:"N/A"`
    Phone   *string `json:"phone" header:"Phone" no:"3" default:"NULL"`
    Email   *string `json:"email" header:"Email" no:"5"`
    Age     *string `json:"age" header:"Age" no:"6" default:"999"`
    Amount  float64 `json:"amount" header:"Amount" no:"7" decimal:"2"`
    Other   string
}
```

## Support type

- `int`, `int8`, `int16`, `int32`, `int64`
- `string`
- `float32`
- `float64`
- pointer values for supported types

## Using for Convert

```go
m := []MyStruct{
    {ID: 1, Name: "N1"},
    {ID: 2, Name: "N2"}
}
csv := csvx.Convert[MyStruct](m)
```

## Using for ConvertWith

`ConvertWith` supports options. Decimal options override the `decimal` tag.

```go
type Product struct {
    ID     int     `json:"id" header:"ID" no:"1"`
    Amount float64 `json:"amount" header:"Amount" no:"2" decimal:"1"`
    Price  float64 `json:"price" header:"Price" no:"3"`
}

products := []Product{
    {ID: 1, Amount: 12.34, Price: 99.999},
}

csv := csvx.ConvertWith(
    products,
    csvx.WithDecimal("Price", 2),
)
```

## ConvertWith options

- `WithDecimal(field, precision)` sets precision for one float field.
- `WithDecimals(map[string]int{...})` sets precision for multiple float fields.
- `WithIgnoreDoubleQuote(true)` disables wrapping cells with double quotes.

The `field` value can be a struct field name, `field` tag, or `header` tag.

## Result

```csv
"ID","Name Space"
"1","N1"
"2","N2"
```

## ConvertWith Result

```csv
"ID","Amount","Price"
"1","12.3","100.00"
```

## Quote every cell (dynamic columns / streaming)

Every cell is wrapped in double quotes and a double quote inside a value is escaped as `""` (RFC 4180),
so values containing `,`, `"` or a line break never shift columns.

```go
csvx.QuoteCell(`say "hi", ok`)            // "say ""hi"", ok"
csvx.FormatQuoted([]string{"1", "a,b"})   // "1","a,b"

// Same API as encoding/csv.Writer, but quotes every cell
w := csvx.NewQuotedWriter(gzipWriter)
_ = w.Write([]string{"ID", "Name"})
_ = w.Write([]string{"1", "Lotus's, Bangna"})
w.Flush()
if err := w.Error(); err != nil {
    // handle
}
```

## Define struct for Parse

Add `header` for mapping in csv header

```go
type Struct struct {
	ID   string `header:"ID"`
	Name string `header:"Name Space"`
}
```

## Using for Parse

```go
rows := [][]string{
    {"ID", "Name Space"},
    {"1", "Name1"},
    {"2", "Name2"},
    {"3", "Name3"},
    {"4", "Name4"},
}
s := csvx.Parser[Struct](rows)
```

## Result

```json
[
  {"ID":"1","Name":"Name1"},
  {"ID":"2","Name":"Name2"},
  {"ID":"3","Name":"Name3"},
  {"ID":"4","Name":"Name4"}
]
```

## Benchmark

```shell
goos: darwin
goarch: arm64
pkg: github.com/prongbang/csvx
cpu: Apple M4 Pro
BenchmarkConvert-12          	  731186	      1609 ns/op	     728 B/op	      35 allocs/op
BenchmarkConvertWith-12      	  695539	      1706 ns/op	     992 B/op	      36 allocs/op
BenchmarkManualConvert-12    	 2286088	      534.1 ns/op	    4346 B/op	       8 allocs/op
BenchmarkTryConvert-12       	 1000000	      1008 ns/op	     496 B/op	      31 allocs/op
```

## 🙏 Acknowledgments

- Built with GO
- IDE Support by [GoLand free 30-day trial.](https://www.jetbrains.com/go/)

---
