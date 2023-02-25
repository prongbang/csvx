# CSVX

Convert array struct to csv format

## Install

```shell
go get github.com/prongbang/csvx
```

## Define struct

Add `field` and `index` start with 1

```go
type MyStruct struct {
    Name string `json:"name" field:"Name Space" index:"2"`
    ID   int    `json:"id" field:"ID" index:"1"`
}
```

## Using

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