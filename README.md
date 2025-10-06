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

Add `header` for mapping in csv header and `no` start with 1 for sort header

```go
type MyStruct struct {
    Name  	string 	`json:"name" header:"Name Space" no:"2"`
    ID    	int    	`json:"id" header:"ID" no:"1"`
    Address *string `json:"address" header:"Address" no:"4" default:"N/A"`
    Phone   *string `json:"phone" header:"Phone" no:"3" default:"NULL"`
    Email   *string `json:"email" header:"Email" no:"5"`
    Age     *string `json:"age" header:"Age" no:"6" default:"999"`
    Other string
}
```

## Support type

- `int`, `int8`, `int16`, `int32`, `int64`
- `string`
- `float64`

## Using for Convert

```go
m := []MyStruct{
    {ID: 1, Name: "N1"},
    {ID: 2, Name: "N2"}
}
csv := csvx.Convert[MyStruct](m)
```

## Result

```csv
"ID","Name Space"
"1","N1"
"2","N2"
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
cpu: Apple M1 Pro
BenchmarkConvert-10          	  430015	      5751 ns/op
BenchmarkManualConvert
BenchmarkManualConvert-10    	 2002738	       614.8 ns/op
BenchmarkTryConvert
BenchmarkTryConvert-10    	      760494	      1573 ns/op
```
