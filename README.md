# CSVX

Convert array struct to csv format and Parse csv format to array struct with Golang 

## Install

```shell
go get github.com/prongbang/csvx
```

## Define struct for Convert

Add `field` for mapping in cav header and `index` start with 1 for sort header

```go
type MyStruct struct {
    Name string `json:"name" field:"Name Space" index:"2"`
    ID   int    `json:"id" field:"ID" index:"1"`
}
```

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

Add `field` for mapping in cav header

```go
type Struct struct {
	ID   string `field:"ID"`
	Name string `field:"Name Space"`
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